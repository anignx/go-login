package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"my.service/go-login/conf"
	"my.service/go-login/package/BackendPlatform/ginhandle"
	myconfig "my.service/go-login/package/config"
	"my.service/go-login/service"
)

var (
	r   *gin.Engine
	svc *service.Service
)

func Init(service *service.Service, conf *conf.Config) {
	svc = service

	gin.SetMode(gin.DebugMode)
	gin.DisableConsoleColor()

	r = gin.Default()
	r.Use(ginhandle.GinLogger())
	r.Use(ginhandle.GinRecovery(true))

	initRouter(r)
	port := fmt.Sprintf(":%d", myconfig.Conf.Server.Port)
	r.Run(port)
}

func Shutdown() {
	if svc != nil {
		svc.Close()
	}
}
