package service

import (
	"github.com/gin-gonic/gin"
	"my.service/go-login/package/BackendPlatform/logging"
)

type UserMessage struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (s *Service) UserLogin(c *gin.Context) (*UserMessage, error) {
	a, _ := s.dao.GetUserByID(1)
	logging.Logger.Infof("user: %v", a)
	return &UserMessage{
		Name: "zhangsan",
		Age:  16,
	}, nil
}

func (s *Service) UserRegister(c *gin.Context) (interface{}, error) {
	// sms.PushIphoneSms("18519121341", "阿里云短信测试", "SMS_154950909", "{\"code\":\"277892\"}")

	//在这里判断各种登录模式（手机验证码登录，第三方账号登录等，由app传过来的参数进行处理），不同的登录模式，调用不同的services函数

	return struct{}{}, nil
}

func (s *Service) UserLogout(c *gin.Context) (interface{}, error) {
	return struct{}{}, nil
}
