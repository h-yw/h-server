package models

import "encoding/xml"

type EventType string
type WXMsgType string

const (
	EventTypeSubscribe   EventType = "subscribe"
	EventTypeUnsubscribe EventType = "unsubscribe"
)
const (
	WXMsgTypeText       WXMsgType = "text"
	WXMsgTypeImage      WXMsgType = "image"
	WXMsgTypeVoice      WXMsgType = "voice"
	WXMsgTypeVideo      WXMsgType = "video"
	WXMsgTypeShortVideo WXMsgType = "shortvideo"
	WXMsgTypeLocation   WXMsgType = "location"
	WXMsgTypeEvent      WXMsgType = "event"
)

type WXReceive struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      WXMsgType
}

type WXTextMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      WXMsgType
	Content      string
	MsgId        int64
}
type WXEventMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Event        string
}

type WXNewsArticle struct {
	Title       string
	Description string
	PicUrl      string
	Url         string
}

type WXNewsReply struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      WXMsgType
	ArticleCount int
	Articles     []WXNewsArticle `xml:"Articles>item,omitempty"`
	XMLName      xml.Name        `xml:"xml"`
}

type WXTextReply struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      WXMsgType
	Content      string
	// 若不标记XMLName, 则解析后的xml名为该结构体的名称
	XMLName xml.Name `xml:"xml"`
}
