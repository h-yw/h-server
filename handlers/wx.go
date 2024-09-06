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
		WXNewsReply(c, rawMsg.ToUserName, rawMsg.FromUserName)
		return
	}
	switch rawMsg.MsgType {
	case models.WXMsgTypeEvent:
		var eventMsg models.WXEventMsg
		if err := xml.Unmarshal(body, &eventMsg); err != nil {
			log.Printf("[消息接收][eventMsg] - 解析文本消息失败: %v\n", err)
			WXNewsReply(c, rawMsg.ToUserName, rawMsg.FromUserName)
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
		if err := xml.Unmarshal(body, &textMsg); err != nil {
			log.Printf("[消息接收][textMsg] - 解析文本消息失败: %v\n", err)
			WXNewsReply(c, rawMsg.ToUserName, rawMsg.FromUserName)
			return
		}
	}
	// err := c.ShouldBindXML(&textMsg)
	// if err != nil {
	// 	log.Printf("[消息接收] - XML数据包解析失败: %v\n", err)
	// }
	// log.Printf("[消息接收] - 收到消息, 消息类型为: %s, 消息内容为: %s\n", textMsg.MsgType, textMsg.Content)

}

func WXNewsReply(c *gin.Context, fromUser, toUser string) {
	replyTextMsg := models.WXTextReply{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      models.WXMsgTypeText,
		Content:      "欢迎来到ifcat！这里将会发布一些技术文章，摄影作品等。当然，你也可以留言，我会回复。\n 你也可以去看我的博客：<a href=\"https://hlovez.life\">hlovez.life</a>",
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
		Content:      fmt.Sprintf("欢迎关注ifcat！这里将会发布一些技术文章，摄影作品等。当然，你也可以留言，我会回复。\n 你也可以去看我的博客：%s", "<a href=\"https://hlovez.life\">hlovez.life</a>"),
	}
	msg, err := xml.Marshal(replyTextMsg)
	if err != nil {
		log.Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)

	}
	_, _ = c.Writer.Write(msg)
}
