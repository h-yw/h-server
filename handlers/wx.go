package handlers

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

func WeixinHandler(c *gin.Context) {
	token := "houyw"
	signature := c.DefaultQuery("signature", "")
	timestamp := c.DefaultQuery("timestamp", "")
	nonce := c.DefaultQuery("nonce", "")
	echostr := c.DefaultQuery("echostr", "")
	if signature == "" || timestamp == "" || nonce == "" || echostr == "" {
		c.String(http.StatusOK, "hello, this is handle view")
	} else {
		list := []string{
			token,
			timestamp,
			nonce,
		}
		sort.Strings(list)
		concatenated := strings.Join(list, "")
		hash := sha1.New()
		hash.Write([]byte(concatenated))
		hashCode := hex.EncodeToString(hash.Sum(nil))
		fmt.Printf("handle/GET func: hashcode, signature: %s, %s\n", hashcode, signature)

		if hashCode == signature {
			c.String(http.StatusOK, echostr)
		} else {
			c.String(http.StatusOK, "")
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
