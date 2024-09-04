// routes/routes.go
package routes

import (
	"hserver/handlers"
	"hserver/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.Use(middlewares.LoggerMiddleware())
	router.Use(middlewares.Cros())
	router.Use(middlewares.Auth())
	v1 := router.Group("/v1")
	{
		v1.GET("/", handlers.HomeHandler)
		v1.GET("/about", handlers.AboutHandler)
		v1.GET("/ping", handlers.PingHandler)
		v1.GET("/ws", handlers.WSHandler)
		v1.POST("/resize", handlers.ResizeHandler)
		v1.POST("/upload", handlers.UploadHandler)
		v1.GET("/puser", handlers.PuserHandler)
		v1.GET("/image.png", handlers.ImageHandler)
	}
	// 添加更多路由和处理函数
}
