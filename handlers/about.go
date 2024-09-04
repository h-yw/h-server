package handlers

import (
	"github.com/gin-gonic/gin"
)

func AboutHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"local": "about",
		"ip":    c.ClientIP(),
	})
}
