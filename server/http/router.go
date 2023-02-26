package http

import "github.com/gin-gonic/gin"

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
}
