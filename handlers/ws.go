package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Msg struct {
	From int32  `json:"from"`
	To   int32  `json:"to"`
	Msg  string `json:"msg"`
}

type Session struct {
	conn *websocket.Conn
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WSHandler(c *gin.Context) {
	// 升级成 websocket 连接
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Fatalln(err)
	}
	// 完成时关闭连接释放资源
	defer conn.Close()
	session := NewSession()
	session.handler(conn)
}

func NewSession() *Session {
	return &Session{}
}
func (s *Session) handler(conn *websocket.Conn) {
	s.conn = conn
	go s.rendLoop()
	// go s.writeLoop()
}
func (s *Session) rendLoop() {
	for {
		// 读取客户端发送过来的消息，如果没发就会一直阻塞住
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			fmt.Println("read error")
			fmt.Println(err)
			break
		}
		fmt.Println(string(message))
		if string(message) == "ping" {
			message = []byte("pong")
		}
		fmt.Println("received msg: %s", message)
	}
}

// func (s *Session) writeLoop() {
// 	for {
// 		select {
// 		case message := <- s.:
// 			err := s.conn.WriteMessage(websocket.TextMessage, []byte(message))
// 			if err != nil {
// 				fmt.Println(err)
// 				break
// 			}
// 		}
// 	}
// }
