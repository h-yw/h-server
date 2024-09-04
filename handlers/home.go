// handlers/home_handler.go
package handlers

import (
	"hserver/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	// c.JSON(200, gin.H{
	// 	"serverName": global.App.Config.App.AppName,
	// 	"email":      "1327603193@qq.com",
	// 	"github":     "https://github.com/h-yw",
	// })
	c.HTML(http.StatusOK, "default.html", gin.H{
		"appName":    global.App.Config.App.AppName,
		"appVersion": global.App.Config.App.Version,
		"author": gin.H{
			"name":    global.App.Config.Author.Name,
			"email":   global.App.Config.Author.Email,
			"twitter": global.App.Config.Author.Twitter,
			"github":  global.App.Config.Author.Github,
		},
	})
}
