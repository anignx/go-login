package daenerys

import (
	"context"
	"fmt"
	"sync"
	"time"

	"my.service/go-login/package/config"
)

const TimeFormat = "2006-01-02 15:04:05.999"

type Daenerys struct {
	initOnce     sync.Once
	mysqlClients sync.Map
	redisClients sync.Map
	ConfigPath   string
	Namespace    string
	App          string
	Name         string
}

func New() *Daenerys {
	return &Daenerys{
		initOnce: sync.Once{},
	}
}

var Defult = New()

func Init(Optional ...Option) error {
	return Defult.Init(Optional...)
}

func (d *Daenerys) Init(options ...Option) error {
	curTime := time.Now().Format(TimeFormat)
	d.initOnce.Do(func() {
		// 把配置所需参数写入daenerys
		for _, option := range options {
			option(d)
		}
		// init middleware client
		if err := d.initMiddleware(); err != nil {

			fmt.Printf("%s init daenerys middleware fatal error:%v, app:%s name:%s namespace:%s config:%s\n",
				curTime, err, d.App, d.Name, d.Namespace, d.ConfigPath)
		}
	})
	return nil
}

func (d *Daenerys) initMiddleware() error {

	// kafkaProducerInit := func() error {
	// 	return d.InitKafkaProducer(d.kafkaProducerConfig())
	// }
	// kafkaConsumerInit := func() error {
	// 	return d.InitKafkaConsume(d.config.KafkaConsume)
	// }
	redisClientInit := func() error {
		return d.InitRedisClient(config.Conf.Redis)
	}
	mysqlClientInit := func() error {
		return d.InitSqlClient(config.Conf.Database)
	}
	// esClientInit := func() error {
	// 	return d.InitESClient(d.config.ESCfg)
	// }

	middlewares := map[string]func() error{
		// "kafkaProducerInit": kafkaProducerInit,
		// "kafkaConsumerInit": kafkaConsumerInit,
		"redisClientInit": redisClientInit,
		"mysqlClientInit": mysqlClientInit,
		// "esClientInit":      esClientInit,
	}
	for name, fn := range middlewares {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		fnDone := make(chan error)
		go func() {
			fnDone <- fn()
		}()
	INNER:
		for {
			select {
			case <-ctx.Done():
				cancel()
				return fmt.Errorf("doing %s timeout, please check your config", name)
			case err := <-fnDone:
				if err != nil {
					cancel()
					return err
				}
				break INNER
			}
		}
		cancel()
	}
	return nil
}
