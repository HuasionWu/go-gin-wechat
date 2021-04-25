// 官方 Demo

package main

import (
	"github.com/gin-gonic/gin"
	"go-gin-weixin/config"
	"go-gin-weixin/pkg/log"
	"go-gin-weixin/router"
)

func main() {
	engine := gin.Default()
	router.InitRouter(engine) // 设置路由
	engine.Run(config.PORT)
	engine.Use(log.LoggerToFile())
}
