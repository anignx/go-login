package util

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	// 手机注册
	MAIN_LOGIN_TYPE = 0
	// 微信注册
	WECHAT_LOGIN_TYPE = 1
	// QQ注册
	QQ_LOGIN_TYPE = 2
	// 微博注册
	WEIBO_LOGIN_TYPE = 3
)

const (
	MID_USER_LOGIN_STRING = "go-login-mid-user-login-string"
)

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
