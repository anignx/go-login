package dao

import (
	"my.service/go-login/conf"
	"my.service/go-login/package/daenerys/proxy"
)

type Dao struct {
	c  *conf.Config
	db *proxy.SQL
}

func New(c *conf.Config) *Dao {
	return &Dao{
		c:  c,
		db: proxy.InitSQL("go-login"),
	}
}

func (d *Dao) Close() error {
	return nil
}
