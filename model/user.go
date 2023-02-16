package model

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordCost        = 12       // 密码加密难度
	Active       string = "active" // 激活用户
)

type User struct {
	gorm.Model
	NickName       string `gorm:"unique"`
	PasswordDigest string
	Email          string
	Avatar         string `gorm:"size:1000"`
	Phone          string
	Status         string
}

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}

// AvatarURL 封面地址
func (user *User) AvatarURL() string {
	return user.Avatar
}
