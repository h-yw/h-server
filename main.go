package main

import (
	"fmt"
	"hserver/bootstrap"
	"hserver/global"
	"hserver/routes"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

type Date struct {
	From int32  `json:"from"`
	To   int32  `json:"to"`
	Msg  string `json:"msg"`
}

func main() {
	// 初始化配置
	bootstrap.InitializeConfig()
	f, _ := os.Create("log/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()
	// staticPath := wd + "/test"
	r.LoadHTMLGlob(global.App.Config.Template.Path)
	r.Static(global.App.Config.StaticServer.Path, "./test")
	routes.SetupRoutes(r)
	url := global.App.Config.Server.Host + ":" + global.App.Config.Server.Port
	fmt.Println("server start at: " + url)
	r.Run(url) // listen and serve on 0.0.0.0:8080
}
