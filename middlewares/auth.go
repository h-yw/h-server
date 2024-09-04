package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JWT struct {
	apiKey string
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// method := c.Request.Method
		// c.Header("Access-Control-Allow-Origin", "*")
		// c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		// c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		// c.Header("Access-Control-Allow-Credentials", "true")
		// if method == "OPTIONS" {
		// 	c.AbortWithStatus(http.StatusNoContent)
		// }
		// FuxUoHyiJwPsb4MR
		var query JWT
		if err := c.ShouldBindQuery(&query); err != nil {
			fmt.Println("query====?", err.Error())
			c.AbortWithStatusJSON(http.StatusNonAuthoritativeInfo, gin.H{
				"msg": "未鉴权",
				"err": err.Error(),
			})
		} else {

			// query := c.GetHeader("Authorization")
			fmt.Println("query====?", query)
			c.Next()
		}

	}
}
