package dao

import "my.service/go-login/package/BackendPlatform/logging"

func (s *Dao) TestRedis() {
	logging.Logger.Info("test redis")
	data, err := s.redis.Do("SET", "test", "test123")
	if err != nil {
		logging.Logger.Infof("set字符串失败，", err)
	}
	logging.Logger.Infof("redis data: %v", data)
}
