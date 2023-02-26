package daenerys

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/clientv3"
	"my.service/go-login/package/config"
	"my.service/go-login/package/daenerys/etcd"
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
	Discovery    *etcd.ServiceDiscovery
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
		d.InitNamespace()
		// 把配置所需参数写入daenerys
		for _, option := range options {
			option(d)
		}
		// init middleware client
		if err := d.initMiddleware(); err != nil {
			fmt.Printf("%s init daenerys middleware fatal error:%v, app:%s name:%s namespace:%s config:%s\n",
				curTime, err, d.App, d.Name, d.Namespace, d.ConfigPath)
		}
		d.ServerRegistry()
	})
	return nil
}

// 服务发现 + 服务注册
func (d *Daenerys) ServerRegistry() error {
	// ip := GetLocalIP()
	port := fmt.Sprintf(":%d", config.Conf.Server.Port)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("etcd boot success")

	hostname, _ := os.Hostname()

	// 创建租约
	svc := &etcd.ServiceRegister{
		Cli: cli,
		Key: fmt.Sprintf("/%s/%s/%s/%s", d.Namespace, d.App, d.Name, hostname),
		Val: "127.0.0.1" + port,
	}
	// 申请租约 5s自动续租
	resp, err := svc.Cli.Grant(context.Background(), 5)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = svc.Cli.Put(context.Background(), svc.Key, svc.Val, clientv3.WithLease(resp.ID))
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 自动续租
	leaseRespChan, err := svc.Cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	svc.LeaseID = resp.ID
	svc.KeepAliveChan = leaseRespChan
	// 处理续租应答
	go func() {
		for {
			select {
			case leaseKeepResp := <-svc.KeepAliveChan:
				if leaseKeepResp == nil {
					fmt.Println("租约失效")
					goto END
				}
			}
		}
	END:
	}()
	fmt.Printf("服务注册成功 key:%s value:%s", svc.Key, svc.Val)
	d.Discovery = etcd.InitDiscovery()
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

func (d *Daenerys) InitNamespace() error {
	d.Namespace = "buzz"
	d.App = config.Conf.Server.App
	d.Name = config.Conf.Server.ServiceName
	return nil
}
