package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"my.service/go-login/dao"
	"my.service/go-login/package/BackendPlatform/logging"
	"my.service/go-login/package/sms"
	"my.service/go-login/util"
)

type RegisterResponse struct {
	UserId uint   `json:"user_id"`
	Msg    string `json:"msg"`
}

type UserLoginResponse struct {
	UserId uint   `json:"user_id"`
	Msg    string `json:"msg"`
}

type SmsCode struct {
	Code string `json:"code"`
}

// user类，暂时先放这里
type User struct {
	UserId   uint   `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func (s *Service) UserLogin(c *gin.Context) (*UserLoginResponse, error) {
	resp := &UserLoginResponse{}
	identifier := c.Query("identifier")
	credential := c.Query("credential")
	identityType, _ := strconv.Atoi(c.Query("identity_type"))
	if identifier == "" || credential == "" {
		return resp, errors.New("identifier or credential is empty")
	}
	logging.Logger.Infof("identifier: %s, credential: %s, identityType: %d", identifier, credential, identityType)
	switch identityType {
	case util.MAIN_LOGIN_TYPE:
		//手机验证码登录，获取redis中的验证码，跟credential进行比对
		code, _ := s.dao.GetLoginCode(identifier)
		if code != credential {
			logging.Logger.Infof("login failed, identifier: %s, credential: %s, code: %v", identifier, credential, code)
			return resp, errors.New("验证码已过期")
		}
		// 验证成功后，删除redis中的验证码
		s.dao.DelLoginCode(identifier)

		// 登录成功
		user, err := s.dao.GetUserByIphone(identifier)
		if err != nil {
			logging.Logger.Infof("login failed, identifier: %s, credential: %s, code: %v", identifier, credential, code)
			return resp, errors.New("登录失败，用户不存在")
		}

		midCode := identifier + util.MID_USER_LOGIN_STRING
		cookieValue := util.Md5(midCode)
		c.SetCookie("login_user", cookieValue, 3600, "/", "localhost", false, true)
		sessionInfo := &SessionInfo{
			UserId:      user.ID,
			Iphone:      user.Iphone,
			CookieValue: cookieValue,
		}

		s.setCurrentUser(c, *sessionInfo)
		logging.Logger.Infof("login success, identifier: %s, credential: %s, code: %v, cookieValue: %v", identifier, credential, code, cookieValue)
		return resp, nil
	default:
		// TODO 第三方登录
		// a, _ := s.dao.GetUserByID(1)
	}

	return resp, errors.New("登录失败")

}

type SessionInfo struct {
	UserId      uint   `json:"user_id"`
	Iphone      string `json:"iphone"`
	CookieValue string `json:"cookie_value"`
}

func (s *Service) getCurrentUser(c *gin.Context) SessionInfo {
	session := sessions.Default(c)
	return session.Get("currentUser").(SessionInfo) // 类型转换一下
}

func (s *Service) setCurrentUser(c *gin.Context, sessionInfo SessionInfo) {
	session := sessions.Default(c)
	session.Set("currentUser", sessionInfo)
	// 一定要Save否则不生效，若未使用gob注册User结构体，调用Save时会返回一个Error
	err := session.Save()
	if err != nil {
		logging.Logger.Errorf("session save error: %v", err)
	}
}

func (s *Service) UserRegister(c *gin.Context) (*RegisterResponse, error) {
	// sms.PushIphoneSms("18519121341", "阿里云短信测试", "SMS_154950909", "{\"code\":\"277892\"}")
	resp := &RegisterResponse{}
	//在这里判断各种登录模式（手机验证码登录，第三方账号登录等，由app传过来的参数进行处理），不同的登录模式，调用不同的services函数
	registerType, _ := strconv.Atoi(c.Query("register_type"))
	logging.Logger.Infof("register_type: %d", registerType)

	var (
		userId uint
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

func (s *Service) MainLoginType(c *gin.Context) (uint, error) {
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

// 获取token时，不进行是否已注册检查，上层调用时可检查
func (s *Service) GetToken(c *gin.Context) (string, error) {
	logging.Logger.Info("get token")
	identifier := c.Query("identifier")
	IdentityType, _ := strconv.Atoi(c.Query("identity_type"))
	if identifier == "" {
		return "", errors.New("identifier is empty")
	}

	switch IdentityType {
	case util.MAIN_LOGIN_TYPE:
		// 获取手机验证码
		ttl, _ := s.dao.CheckLoginCode(identifier)
		if !ttl {
			return "", errors.New("验证码未过期")
		}
		s.dao.DelLoginCode(identifier)
		code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
		jsonCode, _ := json.Marshal(SmsCode{Code: code})

		// TODO 此处需要判断短信是否发送成功
		err := sms.PushIphoneSms(identifier, "阿里云短信测试", "SMS_154950909", string(jsonCode))
		if err != nil {
			logging.Logger.Errorf("push sms error: %v", err)
			return "", err
		}
		// 将code写入redis
		s.dao.SetLoginCode(identifier, code)
	default:
		// TODO 第三方登录
		// a, _ := s.dao.GetUserByID(1)
	}
	return "", nil
}

func (s *Service) UserLogout(c *gin.Context) (interface{}, error) {
	return struct{}{}, nil
}

type UserInfoResponse struct {
	UserId   uint   `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Iphone   string `json:"iphone"`
}

// 登录时，获取用户信息，非登录时不可用
func (s *Service) GetUserInfo(c *gin.Context) (*UserInfoResponse, error) {
	sessionInfo := s.getCurrentUser(c)
	logging.Logger.Infof("sessionInfo: %v", sessionInfo)
	cookieValue, _ := c.Cookie("login_user")
	if cookieValue != sessionInfo.CookieValue {
		logging.Logger.Infof("cookieValue: %v, sessionInfo: %v", cookieValue, sessionInfo)
		return nil, errors.New("个人信息获取失败，请重新登录")
	}

	resp := &UserInfoResponse{}

	userId := sessionInfo.UserId
	if userId == 0 {
		return resp, errors.New("user_id is empty")
	}
	user, err := s.dao.GetUserByID(userId)
	if err != nil {
		logging.Logger.Errorf("get user error: %v", err)
		return resp, err
	}
	resp.UserId = user.ID
	resp.Nickname = user.Nickname
	resp.Avatar = user.Avatar
	resp.Iphone = user.Iphone
	return resp, nil
}
