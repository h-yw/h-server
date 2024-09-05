package handlers

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"hserver/models"
	"log"
	"net/http"
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
	var textMsg models.WXTextMsg
	err := c.ShouldBindXML(&textMsg)
	if err != nil {
		log.Printf("[消息接收] - XML数据包解析失败: %v\n", err)
	}
	log.Printf("[消息接收] - 收到消息, 消息类型为: %s, 消息内容为: %s\n", textMsg.MsgType, textMsg.Content)
	WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName)
}

func WXMsgReply(c *gin.Context, fromUser, toUser string) {
	replyTextMsg := models.WXNewsReply{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      "news",
		ArticleCount: 2,
		Articles: []models.WXNewsArticle{
			{
				Title:       "欢迎访问博客",
				Description: "欢迎访问博客",
				PicUrl:      "https://hlovez.life/static/favicons/logo_800x320.png",
				Url:         "https://hlovez.life",
			},
			{
				Title:       "使用 GitHub Actions 自动构建和部署 Docker 镜像",
				Description: "欢迎访问博客",
				PicUrl:      "https://mmbiz.qpic.cn/sz_mmbiz_jpg/bQibWia0UpTvJqJIJlwUhibsB04k9sg27ZhXVahPdW3O4HDASrGzPVSoxzYfJZSMAibsZ3cNuhGo3X2EkXMVdU85rQ/0?wx_fmt=jpeg",
				Url:         "https://mp.weixin.qq.com/s/Y6mQZD9VxZNMWznQbeoFxA?token=301568447&lang=zh_CN",
			},
		},
	}
	msg, err := xml.Marshal(replyTextMsg)
	if err != nil {
		log.Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)

	}
	_, _ = c.Writer.Write(msg)
}
