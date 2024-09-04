// middlewares/logger.go
package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("请求来自: " + c.ClientIP())
		c.Next()
	}
}
