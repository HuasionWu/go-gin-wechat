package router

import (
	"github.com/gin-gonic/gin"
	"go-gin-weixin/wechat"
)

func InitRouter(r *gin.Engine) {

	GroupWx := r.Group("/wx")
	{
		//微信公众号服务器配置接口
		//GroupWx.GET("/message", wechat.ServeHTTP)
		//微信公众号自定义菜单接口
		GroupWx.GET("/menu", wechat.CreateMenu)
		//生成带参数的二维码（微信公众号）
		GroupWx.GET("/code", wechat.GetCode)
		//微信公众号消息接收
		GroupWx.POST("/message", wechat.WXMsgReceive)
	}
}
