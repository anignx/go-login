package service

import (
	"my.service/go-login/conf"
	"my.service/go-login/dao"
)

type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

func New(conf *conf.Config) *Service {
	return &Service{
		c:   conf,
		dao: dao.New(conf),
	}
}

func (s *Service) Close() {
	// 关闭数据库连接
	// 关闭redis连接
}
