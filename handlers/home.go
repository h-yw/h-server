// handlers/home_handler.go
package handlers

import (
	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"title": "Home",
		"ip":    c.ClientIP(),
	})
}
