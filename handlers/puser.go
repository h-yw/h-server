package handlers

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func PuserHandler(c *gin.Context) {
	path, _ := os.Getwd()
	if pusher := c.Writer.Pusher(); pusher != nil {
		// 使用 pusher.Push() 做服务器推送
		if err := pusher.Push(path+"/test/app.js", nil); err != nil {
			fmt.Printf("Failed to push: %v", err)
		}
	}
	c.HTML(200, "index.html", gin.H{
		"status": "success",
	})
}
