package service

import (
	"encoding/gob"

	"my.service/go-login/conf"
	"my.service/go-login/dao"
	"my.service/go-login/manager"
)

type Service struct {
	c   *conf.Config
	dao *dao.Dao
	mgr *manager.Manager
}

func New(conf *conf.Config) *Service {
	s := &Service{
		c:   conf,
		dao: dao.New(conf),
		mgr: manager.New(conf),
		//kafka， 定时任务等
	}

	// cookie注册
	gob.Register(SessionInfo{})
	return s
}

func (s *Service) Close() {
	// 关闭数据库连接
	// 关闭redis连接
}
