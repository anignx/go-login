package service

import (
	"encoding/gob"

	"my.service/go-login/conf"
	"my.service/go-login/dao"
)

type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

func New(conf *conf.Config) *Service {
	s := &Service{
		c:   conf,
		dao: dao.New(conf),
	}

	gob.Register(SessionInfo{})
	return s
}

func (s *Service) Close() {
	// 关闭数据库连接
	// 关闭redis连接
}
