package request

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"
)

type requestCreate struct {
	BaseURL  string
	Client   *http.Client
	Timeout  int
	attempts int
}

type GetParams struct {
	Url     string
	Query   map[string]string
	Headers map[string]string
}
type PostParams struct {
	Url     string
	Body    string
	Query   map[string]string
	Headers map[string]string
}
type RequestResData struct {
}
type RequestRes struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

func (r *requestCreate) Get(params GetParams) (*RequestRes, error) {
	// 拼接BaseUrl和url
	params.Url, _ = MergeURL(r.BaseURL, params.Url)
	log.Printf("url:%s", params.Url)
	req, err := http.NewRequest("GET", params.Url, nil)
	if err != nil {
		log.Println("creating request: %v", err)
		return nil, err
	}
	q := req.URL.Query()
	if params.Query != nil {
		for k, v := range params.Query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if params.Headers != nil {
		for k, v := range params.Headers {
			req.Header.Add(k, v)
		}
	}
	resp, err := r.Client.Do(req)
	if err != nil {
		log.Printf("sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read body error: %v", err)
		return nil, err
	}

	var data map[string]interface{}
	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		log.Printf("unmarshal body error: %v", jsonErr)
		return nil, jsonErr
	}
	result := RequestRes{
		Status: resp.Status,
		Data:   data,
	}

	return &result, nil
}
func (r *requestCreate) Post(params PostParams) (*RequestRes, error) {
	// 拼接BaseUrl和url
	params.Url, _ = MergeURL(r.BaseURL, params.Url)
	log.Printf("url:%s", params.Url)
	paramsUrl, _ := url.Parse(params.Url)

	q := paramsUrl.Query()
	if params.Query != nil {
		for k, v := range params.Query {
			q.Add(k, v)
		}
	}
	q.Add("access_token", accessToken)
	paramsUrl.RawQuery = q.Encode()

	req, err := http.NewRequest("POST", paramsUrl.String(), bytes.NewBufferString(params.Body))
	if err != nil {
		log.Printf("creating request: %v", err)
		return nil, err
	}
	if params.Headers != nil {
		for k, v := range params.Headers {
			req.Header.Add(k, v)
		}
	}
	resp, err := r.Client.Do(req)
	if err != nil {
		log.Printf("sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read body error: %v", err)
		return nil, err
	}
	var data map[string]interface{}
	log.Printf("body:%s", string(body))
	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		log.Printf("unmarshal body error: %v", jsonErr)
		return nil, jsonErr
	}
	result := RequestRes{
		Status: resp.Status,
		Data:   data,
	}
	return &result, nil
}

func NewRequest(baseUrl string) *requestCreate {
	GetAccessToken()
	return &requestCreate{
		BaseURL: "https://api.weixin.qq.com/cgi-bin/",
		Client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		attempts: 3,
	}
}

var (
	accessToken string
	expiresIn   time.Time
	mu          sync.Mutex
)

func GetAccessToken() {
	if len(accessToken) > 0 || !time.Now().After(expiresIn) {
		return
	}
	mu.Lock()
	accessToken, expiresIn = getAccessToken()
	mu.Unlock()
}

// 直接使用NewRequest引起循环引用，
func getAccessToken() (string, time.Time) {
	client := &http.Client{}

	reqUrl, err := url.Parse("https://api.weixin.qq.com/cgi-bin/token")
	if err != nil {
		log.Printf("url parse error:%v", err)
		return "", time.Now()
	}
	q := reqUrl.Query()
	q.Add("grant_type", "client_credential")
	q.Add("appid", "wx5c6aa6d21e4e6a40")
	q.Add("secret", "d5caec27c9ce075a5dbccc2359e76e18")
	reqUrl.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", reqUrl.String(), nil)

	res, err := client.Do(req)
	if err != nil {
		log.Printf("request error:%v", err)
		return "", time.Now()
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("read body error:%v", err)
		return "", time.Now()
	}
	var data map[string]interface{}
	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		log.Printf("unmarshal body error:%v", jsonErr)
		return "", time.Now()
	}

	if val, exists := data["access_token"]; exists {
		expires, _ := data["expires_in"].(int)

		return val.(string), time.Now().Add(time.Duration(expires) * time.Second)
	}
	// 获取 access_token
	return "", time.Now()
}

func MergeURL(baseURL, pathname string) (string, error) {
	relative, err := url.Parse(pathname)
	if err != nil {
		return "", err
	}

	// 解析 relativeURL

	if err != nil {
		return "", err
	}

	if relative.IsAbs() {
		return pathname, nil
	}

	// 解析 baseURL
	base, err := url.Parse(baseURL)

	if err != nil {
		return "", err
	}

	// 处理 baseURL 和 relativeURL 的路径合并
	base.Path = path.Join(base.Path, relative.Path)

	return base.String(), nil
}
