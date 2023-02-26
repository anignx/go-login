package dao

import (
	"github.com/garyburd/redigo/redis"
	"my.service/go-login/package/BackendPlatform/logging"
)

func (s *Dao) TestRedis() {
	logging.Logger.Info("test redis")
	data, err := s.redis.Do("SET", "test", "test123")
	if err != nil {
		logging.Logger.Infof("set字符串失败，", err)
	}
	logging.Logger.Infof("redis data: %v", data)
}

func (s *Dao) CheckLoginCode(identifier string) (bool, error) {
	ttl, _ := redis.Int(s.loginCode.Do("TTL", identifier))
	if ttl > 240 {
		return false, nil
	}
	return true, nil
}

func (s *Dao) HasLoginCode(identifier string) bool {
	ttl, _ := redis.Int(s.loginCode.Do("TTL", identifier))
	return ttl > 0
}

func (s *Dao) GetLoginCode(identifier string) (string, error) {
	if !s.HasLoginCode(identifier) {
		// code已失效
		return "", nil
	}
	code, err := redis.String(s.loginCode.Do("GET", identifier))
	if err != nil {
		return "", err
	}
	return code, nil
}

func (s *Dao) DelLoginCode(identifier string) error {
	s.loginCode.Do("DEL", identifier)
	return nil
}

func (s *Dao) SetLoginCode(identifier string, code string) error {
	s.loginCode.Do("SET", identifier, code)
	s.loginCode.Do("EXPIRE", identifier, 300)
	logging.Logger.Infof("set login %s code: %s", identifier, code)
	return nil
}
