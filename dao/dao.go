package dao

import (
	"github.com/gomodule/redigo/redis"
	"my.service/go-login/conf"
	"my.service/go-login/package/daenerys/proxy"
)

type Dao struct {
	c     *conf.Config
	db    *proxy.SQL
	redis redis.Conn
}

func New(c *conf.Config) *Dao {
	return &Dao{
		c:     c,
		db:    proxy.InitSQL("go-login"),
		redis: proxy.RedisClient("login"),
	}
}

func (d *Dao) Close() error {
	return nil
}
