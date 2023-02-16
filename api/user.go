package api

import (
	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
	"github.com/theone-daxia/chat-demo/service"
)

func UserRegister(c *gin.Context) {
	var us service.UserRegisterService
	if err := c.ShouldBind(&us); err != nil {
		logging.Info(err)
		c.JSON(400, ErrorResponse(err))
	} else {
		res := us.Register()
		c.JSON(200, res)
	}
}
