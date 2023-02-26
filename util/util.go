package util

import (
	"crypto/md5"
	"encoding/hex"
	"net"
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

// GetLocalIP 获取本机IP
func GetLocalIP() ([]string, error) {
	ret := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ret, err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ret = append(ret, ipnet.IP.String())
			}
		}
	}
	return ret, err
}
