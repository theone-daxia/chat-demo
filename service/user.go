package service

import (
	logging "github.com/sirupsen/logrus"
	"github.com/theone-daxia/chat-demo/model"
	"github.com/theone-daxia/chat-demo/pkg/e"
	"github.com/theone-daxia/chat-demo/serializer"
)

type UserRegisterService struct {
	NickName string `form:"nick_name" json:"nick_name" binding:"required,min=1,max=10"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=10"`
}

func (us *UserRegisterService) Register() serializer.Response {
	var count int64
	code := e.Success

	model.DB.Model(&model.User{}).Where("nick_name=?", us.NickName).Count(&count)
	if count > 0 {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "该用户名已存在",
		}
	}

	user := model.User{
		NickName: us.NickName,
		Status:   model.Active,
		Avatar:   "https://q1.qlogo.cn/g?b=qq&nk=294350394&s=640",
	}

	// 设置密码（加密）
	if err := user.SetPassword(us.Password); err != nil {
		logging.Info(err)
		code := e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	// 创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
