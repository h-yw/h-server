package handlers

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"hserver/models"
	"hserver/request"
	"hserver/utils"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func WXCheckSign(c *gin.Context) {
	token := "houyw"
	signature := c.DefaultQuery("signature", "")
	timestamp := c.DefaultQuery("timestamp", "")
	nonce := c.DefaultQuery("nonce", "")
	echostr := c.DefaultQuery("echostr", "")
	if signature == "" || timestamp == "" || nonce == "" || echostr == "" {
		c.String(http.StatusOK, "hello, this is handle view")
	} else {
		list := []string{
			token,
			timestamp,
			nonce,
		}
		sort.Strings(list)
		concatenated := strings.Join(list, "")
		hash := sha1.New()
		hash.Write([]byte(concatenated))
		hashCode := hex.EncodeToString(hash.Sum(nil))
		fmt.Printf("handle/GET func: hashcode, signature: %s, %s\n", hashCode, signature)

		if hashCode == signature {
			c.String(http.StatusOK, echostr)
		} else {
			c.String(http.StatusOK, "")
		}
	}
}
func WXMsgReceive(c *gin.Context) {
	// c.ShouldBindXML() 只能读一次
	// 获取原始请求体
	body, err := c.GetRawData()
	if err != nil {
		log.Printf("[消息接收] - 获取原始数据失败: %v\n", err)
		return
	}
	var rawMsg models.WXReceive
	if err := xml.Unmarshal(body, &rawMsg); err != nil {
		log.Printf("[消息接收][rawMsg] - XML数据包解析失败: %v\n", err)
		WXNewsReply(c, rawMsg.ToUserName, rawMsg.FromUserName, "")
		return
	}
	switch rawMsg.MsgType {
	case models.WXMsgTypeEvent:
		var eventMsg models.WXEventMsg
		if err := xml.Unmarshal(body, &eventMsg); err != nil {
			log.Printf("[消息接收][eventMsg] - 解析文本消息失败: %v\n", err)
			WXNewsReply(c, rawMsg.ToUserName, rawMsg.FromUserName, "")
			return
		}
		log.Printf("[消息接收] - 接收到事件消息: %v\n%v\n", eventMsg.Event, models.EventTypeSubscribe)
		if eventMsg.Event == string(models.EventTypeSubscribe) {
			log.Printf("[消息接收] - 接收到关注事件消息: %v\n", eventMsg.Event)
			WXSubscribeReply(c, rawMsg.ToUserName, rawMsg.FromUserName)
		}
		return
	case models.WXMsgTypeText:
		var textMsg models.WXTextMsg
		var str string = ""
		if err := xml.Unmarshal(body, &textMsg); err != nil {
			log.Printf("[消息接收][textMsg] - 解析文本消息失败: %v\n", err)
		}
		feature, err := splitText(textMsg.Content)
		if err != nil {
			log.Printf("[消息接收] - 解析feature失败: %v\n", err)
		} else {
			str = featureHandle(feature)
		}

		WXNewsReply(c, rawMsg.ToUserName, rawMsg.FromUserName, str)
	}
	// err := c.ShouldBindXML(&textMsg)
	// if err != nil {
	// 	log.Printf("[消息接收] - XML数据包解析失败: %v\n", err)
	// }
	// log.Printf("[消息接收] - 收到消息, 消息类型为: %s, 消息内容为: %s\n", textMsg.MsgType, textMsg.Content)

}

func WXNewsReply(c *gin.Context, fromUser, toUser string, content string) {
	defaultStr := "🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉\n欢迎来到ifcat🐱！这里将会发布一些技术文章，摄影作品等。当然，你也可以留言，我会回复😁。\n你也可以去看我的博客💻<a href=\"https://hlovez.life\">hlovez.life</a>\n\n功能列表：\n\t\t<a href=\"#\" style=\"color:#167829\">翻译：</a>\n\t\t\t\t输入例子：\n\t\t\t\t\t\t[trans]这是要翻译的内容"
	if content != "" {
		defaultStr = content
	}
	replyTextMsg := models.WXTextReply{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      models.WXMsgTypeText,
		Content:      defaultStr,
	}
	msg, err := xml.Marshal(replyTextMsg)
	if err != nil {
		log.Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)

	}
	_, _ = c.Writer.Write(msg)
}

func WXSubscribeReply(c *gin.Context, fromUser, toUser string) {

	replyTextMsg := models.WXTextReply{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      models.WXMsgTypeText,
		Content:      fmt.Sprintf("🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉\n欢迎关注ifcat🐱！这里将会发布一些技术文章，摄影作品等。\n当然，你也可以留言，我会回复😁。\n你也可以去看我的博客💻%s\n功能列表：\n\t\t<a href=\"#\" style=\"color:#167829\">翻译</a>：\n\t\t\t\t输入例子：\n\t\t\t\t\t\t[trans]这是要翻译的内容", "<a href=\"https://hlovez.life\">hlovez.life</a>"),
	}
	msg, err := xml.Marshal(replyTextMsg)
	if err != nil {
		log.Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)

	}
	_, _ = c.Writer.Write(msg)
}

func translate(content string) *string {
	var trans string
	req := request.NewRequest("https://api.weixin.qq.com/cgi-bin/")
	lang := utils.AutoCheckLang(&content)
	log.Printf("[translate] - 翻译前: %v\n", lang)
	res, _ := req.Post(request.PostParams{
		Url: "media/voice/translatecontent",
		Query: map[string]string{
			// "access_token": token,
			"lfrom": lang.From,
			"lto":   lang.To,
		},
		Body: content,
	})
	log.Printf("[translate] - 翻译后: %v\n", res.Data)
	if val, exits := res.Data["to_content"]; exits {
		trans = val.(string)
	}
	return &trans
}
func splitText(str string) (*models.Feature, error) {
	re := regexp.MustCompile(`(?s)(\[.*?\])(.*)`)
	matches := re.FindStringSubmatch(str)
	if len(matches) == 0 {
		return nil, fmt.Errorf("未匹配到featureFlag\n")
	}
	var feature models.Feature
	log.Printf("[splitText] - featureFlag: %v, featureValue: %v\n", matches[1], matches[2])
	if val := matches[1]; len(val) > 0 {
		feature.Flag = val
	}
	if val := matches[2]; len(val) > 0 {
		feature.Value = val
	}
	return &feature, nil
}

// 实现feature功能
func featureHandle(feature *models.Feature) string {
	var content string
	switch feature.Flag {
	case string(models.FlagTrans):
		val := translate(feature.Value.(string))
		content = fmt.Sprintf("[翻译结果]\n\n %s", *val)
	default:
		log.Printf("[featureHandle] - 未定义的feature: %v\n", feature.Flag)
		content = ""
	}
	return content
}
