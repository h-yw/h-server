package handlers

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg" // 通过 jpeg 包中的 init 函数注册解码器
	_ "image/png"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/draw"
)

func ResizeHandler(c *gin.Context) {

	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Open(path + "/test/hg.jpg")
	output, err := os.Create(path + "/test/newhg.jpg")
	defer f.Close()
	defer output.Close()
	if err != nil {
		fmt.Println("f=====>", err)
	}
	img, formatName, err := image.Decode(f) // image.Decode(f)
	if err != nil {
		fmt.Println("img====>", err)
	}
	origin_width := img.Bounds().Dx()
	origin_height := img.Bounds().Dy()

	dst := image.NewRGBA(image.Rect(0, 0, origin_width*2, origin_height*2))
	draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

	options := &jpeg.Options{
		Quality: 100,
	}
	jpeg.Encode(output, dst, options)
	c.JSON(200, gin.H{
		"type":       "IMAGE_RESIZE",
		"formatName": formatName,
		"size": gin.H{
			"width":  origin_width,
			"height": origin_height,
		},
		"image": "http://127.0.0.1:9527/static/test/newhg.jpg",
	})
}
