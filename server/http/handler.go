package http

import (
	"github.com/gin-gonic/gin"
	"my.service/go-login/package/BackendPlatform/gec"
	"my.service/go-login/package/BackendPlatform/logging"
)

func login(c *gin.Context) {
	data, err := svc.UserLogin(c)
	if err != nil {
		logging.Logger.Infof("login failed: %v", err)
		c.JSON(200, gin.H{
			"message": "login failed",
		})
		return
	}
	// gec.JSON(c, nil, gec.ErrBadRequest)
	gec.Success(data).JSON(c)
}

func register(c *gin.Context) {
	data, err := svc.UserRegister(c)
	if err != nil {
		logging.Logger.Infof("register failed: %v", err)
		c.JSON(200, gin.H{
			"message": "register failed",
		})
		return
	}
	gec.Success(data).JSON(c)
}

func gettoken(c *gin.Context) {
	data, err := svc.GetToken(c)
	if err != nil {
		logging.Logger.Infof("gettoken failed: %v", err)
		c.JSON(200, gin.H{
			"message": "gettoken failed",
		})
		return
	}
	gec.Success(data).JSON(c)
}

func logout(c *gin.Context) {
	data, err := svc.UserLogout(c)
	if err != nil {
		logging.Logger.Infof("logout failed: %v", err)
		c.JSON(200, gin.H{
			"message": "logout failed",
		})
		return
	}
	gec.Success(data).JSON(c)
}

func getuserinfo(c *gin.Context) {
	data, err := svc.GetUserInfo(c)
	if err != nil {
		logging.Logger.Infof("getuserinfo failed: %v", err)
		c.JSON(200, gin.H{
			"message": "getuserinfo failed",
		})
		return
	}
	gec.Success(data).JSON(c)
}
