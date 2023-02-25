package service

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"my.service/go-login/dao"
	"my.service/go-login/package/BackendPlatform/logging"
	"my.service/go-login/util"
)

type UserMessage struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type RegisterResponse struct {
	UserId int    `json:"user_id"`
	Msg    string `json:"msg"`
}

func (s *Service) UserLogin(c *gin.Context) (*UserMessage, error) {
	a, _ := s.dao.GetUserByID(1)
	logging.Logger.Infof("user: %v", a)
	return &UserMessage{
		Name: "zhangsan",
		Age:  16,
	}, nil
}

func (s *Service) UserRegister(c *gin.Context) (*RegisterResponse, error) {
	// sms.PushIphoneSms("18519121341", "阿里云短信测试", "SMS_154950909", "{\"code\":\"277892\"}")
	resp := &RegisterResponse{}
	//在这里判断各种登录模式（手机验证码登录，第三方账号登录等，由app传过来的参数进行处理），不同的登录模式，调用不同的services函数
	registerType, _ := strconv.Atoi(c.Request.FormValue("register_type"))
	logging.Logger.Infof("register_type: %d", registerType)

	var (
		userId int
		err    error
	)

	switch registerType {
	case util.MAIN_LOGIN_TYPE:
		//主登录方式
		userId, err = s.MainLoginType(c)
	case util.WECHAT_LOGIN_TYPE:
	case util.QQ_LOGIN_TYPE:
	case util.WEIBO_LOGIN_TYPE:
	default:
	}
	if err != nil || userId == 0 {
		logging.Logger.Errorf("register error userId empty: %v", err)
		return resp, err
	}
	resp.UserId = userId
	resp.Msg = "注册成功"
	// 注册成功后续操作
	logging.Logger.Infof("register success, user id: %d", userId)
	return resp, nil
}

func (s *Service) UserLogout(c *gin.Context) (interface{}, error) {
	return struct{}{}, nil
}

func (s *Service) MainLoginType(c *gin.Context) (int, error) {
	nickname := c.Query("nickname")
	avatar := c.Query("avatar")
	iphone := c.Query("iphone")
	credential := c.Query("credential")
	logging.Logger.Debugf("nickname: %s, avatar: %s, iphone: %s", nickname, avatar, iphone)
	if nickname == "" || avatar == "" || iphone == "" {
		return 0, errors.New("params is empty")
	}
	// 调用dao层的函数
	user := &dao.User{
		Nickname: nickname,
		Iphone:   iphone,
		Avatar:   avatar,
	}
	userId, err := s.dao.CreateUser(user, util.MAIN_LOGIN_TYPE, credential)
	if err != nil {
		logging.Logger.Errorf("create user error: %v", err)
		return 0, err
	}
	return userId, nil
}
