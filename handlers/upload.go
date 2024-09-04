package handlers

import (
	"hserver/global"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
	path, _ := os.Getwd()
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	url := global.App.Config.Server.Protocol + "://" + global.App.Config.Server.Host + ":" + global.App.Config.Server.Port + global.App.Config.StaticServer.Path + "/"

	var err error
	var urls []string
	for _, file := range files {
		dst := path + "/test/" + file.Filename
		err = c.SaveUploadedFile(file, dst)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 500,
				"msg":  "上传失败: " + err.Error(),
				"data": gin.H{
					"filename": file.Filename,
				},
			})
			break
		}

		urls = append(urls, url+file.Filename)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "上传成功",
		"data": gin.H{
			"urls": urls,
		},
	})
}
