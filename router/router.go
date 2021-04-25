package router

import (
	"github.com/gin-gonic/gin"
	"go-gin-weixin/wechat"
)

func InitRouter(r *gin.Engine) {

	GroupWx := r.Group("/wx")
	{
		//微信公众号服务器配置接口
		GroupWx.GET("/message", wechat.ServeHTTP)
		//微信公众号自定义菜单接口
		GroupWx.GET("/menu", wechat.CreateMenu)

		//微信公众号二维码生成接口
		GroupWx.GET("/code", wechat.GetCode)
	}
}
