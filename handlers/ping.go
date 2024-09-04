package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"handlerName": "ping",
		"ip":          c.Request.RemoteAddr,
	})
}
