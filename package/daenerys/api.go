package daenerys

import (
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"my.service/go-login/package/BackendPlatform/logging"
	"my.service/go-login/package/config"
)

func (d *Daenerys) InitSqlClient(sqlList []config.SQLGroupConfig) error {
	configs := make([]interface{}, len(sqlList))
	for index, defaultConf := range sqlList {
		configs[index] = defaultConf
	}

	for _, conf := range configs {
		var err error
		sqlConf, ok := conf.(config.SQLGroupConfig)
		if !ok {
			fmt.Printf("sql init err, conf: %v", conf)
			continue
		}
		if _, ok := d.mysqlClients.Load(sqlConf.Name); ok {
			continue
		}

		var g *gorm.DB
		if len(sqlConf.LogLevel) == 0 {
			sqlConf.LogLevel = strings.ToLower(config.Conf.Log.Level)
		}
		// 先只初始化主库，从库再说
		g, err = OpenMySQLDB(sqlConf.Master)
		if err != nil {
			fmt.Printf("sql init err, name: %s, err: %s", sqlConf.Name, err.Error())
			continue
		}
		d.mysqlClients.LoadOrStore(sqlConf.Name, g)
	}
	return nil
}

func (d *Daenerys) SQLClient(name string) *gorm.DB {
	if client, ok := d.mysqlClients.Load(name); ok {
		if v, ok1 := client.(*gorm.DB); ok1 {
			return v
		}
	}
	fmt.Printf("namespace %s mysql client for %s not exist\n", d.Namespace, name)
	logging.Logger.Errorf("namespace %s mysql client for %s not exist", d.Namespace, name)
	return nil
}

func OpenMySQLDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}

func (d *Daenerys) RedisClient(name string) redis.Conn {
	if client, ok := d.redisClients.Load(name); ok {
		if v, ok1 := client.(*redis.Pool); ok1 {
			return v.Get()
		}
	}

	fmt.Printf("namespace %s redis client for %s not exist\n", d.Namespace, name)
	logging.Logger.Infof("namespace %s redis client for %s not exist", d.Namespace, name)
	return nil
}

func (d *Daenerys) InitRedisClient(rcList []config.RedisConfig) error {
	configs := make([]interface{}, len(rcList))
	for index, defaultConf := range rcList {
		configs[index] = defaultConf
	}

	for _, conf := range configs {
		var err error

		redisConf, ok := conf.(config.RedisConfig)
		if !ok {
			fmt.Printf("redis conf type err")
			continue
		}
		if _, ok := d.redisClients.Load(redisConf.ServerName); ok {
			continue
		}

		var client *redis.Pool
		cc := redisConf
		client, err = OpenRedis(&cc)
		if err != nil {
			fmt.Printf("redis init err, serverName: %s, err: %s", redisConf.ServerName, err.Error())
			continue
		}
		d.redisClients.LoadOrStore(cc.ServerName, client)
	}
	fmt.Printf("redis init success, namespace: %s, %s", d.Namespace, configs[0].(config.RedisConfig).Addr)

	return nil
}

// 这里返回redis.pool
func OpenRedis(conf *config.RedisConfig) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: time.Duration(conf.IdleTimeout),
		Dial: func() (redis.Conn, error) {
			pass := redis.DialPassword(conf.Password)
			intRT := time.Duration(conf.ReadTimeout) * time.Millisecond
			readTimeout := redis.DialReadTimeout(intRT)
			intWT := time.Duration(conf.WriteTimeout) * time.Millisecond
			writeTimeout := redis.DialWriteTimeout(intWT)
			intCT := time.Duration(conf.ConnectTimeout) * time.Millisecond
			connTimeout := redis.DialConnectTimeout(intCT)
			databases := redis.DialDatabase(conf.Database)
			var (
				c   redis.Conn
				err error
			)
			c, err = redis.Dial("tcp", conf.Addr, pass, readTimeout, writeTimeout, connTimeout, databases)

			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	return pool, nil
}
