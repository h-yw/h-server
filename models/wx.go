package models

import "encoding/xml"

type WXTextMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	MsgId        int64
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
	MsgType      string
	ArticleCount int
	Articles     []WXNewsArticle
	XMLName      xml.Name `xml:"xml"`
}

type WXTextReply struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	// 若不标记XMLName, 则解析后的xml名为该结构体的名称
	XMLName xml.Name `xml:"xml"`
}
