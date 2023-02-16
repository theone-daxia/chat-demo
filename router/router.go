package router

import (
	"github.com/gin-gonic/gin"
	"github.com/theone-daxia/chat-demo/api"
	"github.com/theone-daxia/chat-demo/service"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery(), gin.Logger())

	v1 := r.Group("/")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "success")
		})
		v1.POST("user/register", api.UserRegister)
		v1.GET("ws", service.WsHandler)
	}
	return r
}
