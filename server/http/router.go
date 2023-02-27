package http

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

func initRouter(s *gin.Engine) {
	s.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	s.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "api",
		})
	})
	s.GET("/login", login)
	s.GET("/register", register)
	s.GET("/logout", logout)
	s.GET("/gettoken", gettoken)
	s.GET("/getuserinfo", getuserinfo)

	ServerConsul(r)
}

func ServerConsul(s *gin.Engine) error {
	// port := fmt.Sprintf(":%d", config.Conf.Server.Port)
	consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))

	// server := web.NewService(
	// 	web.Name("ProductService"), // 当前微服务服务名
	// 	web.Registry(cr),           // 注册到consul
	// 	web.Address(port),          // 端口
	// 	web.Metadata(map[string]string{"protocol": "http"}), // 元信息
	// 	web.Handler(s)) // 路由
	// _ = server.Init()
	// _ = server.Run()
	return nil
}
