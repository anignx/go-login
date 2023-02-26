package etcd

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
	clientv3 "go.etcd.io/etcd/clientv3"
	"my.service/go-login/package/config"
)

//ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	cli        *clientv3.Client  //etcd client
	serverList map[string]string //服务列表
	lock       sync.Mutex
	myService  map[string]struct{} //本服务名称
}

//NewServiceDiscovery  新建发现服务
func NewServiceDiscovery(endpoints []string) *ServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &ServiceDiscovery{
		cli:        cli,
		serverList: make(map[string]string),
	}
}

//WatchService 初始化服务列表和监视
func (s *ServiceDiscovery) WatchService(prefix string) error {
	//根据前缀获取现有的key
	resp, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		s.SetServiceList(string(ev.Key), string(ev.Value))
	}

	//监视前缀，修改变更的server
	watchStartRevision := resp.Header.Revision + 1
	go s.watcher(prefix, watchStartRevision)
	return nil
}

//watcher 监听前缀
func (s *ServiceDiscovery) watcher(prefix string, watchStartRevision int64) {
	rch := s.cli.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithRev(watchStartRevision))
	log.Printf("watching prefix:%s watchStartRevision:%v, now...", prefix, watchStartRevision)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT: //修改或者新增
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: //删除
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

//SetServiceList 新增服务地址
func (s *ServiceDiscovery) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[key] = string(val)
	log.Println("put key :", key, "val:", val)
}

//DelServiceList 删除服务地址
func (s *ServiceDiscovery) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, key)
	log.Println("del key:", key)
}

//GetServices 获取服务地址
func (s *ServiceDiscovery) GetServices() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]string, 0)

	// TODO 需要获取所有服务的列表，写入到map中，然后根据map的key来获取服务
	for _, v := range s.serverList {
		addrs = append(addrs, v)
	}
	return addrs
}

//Close 关闭服务
func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}

func (s *ServiceDiscovery) InitMyService() {
	s.myService = make(map[string]struct{})
	for _, v := range config.Conf.ServerClient {
		s.myService[v.ServiceName] = struct{}{}
	}
}

func InitDiscovery() *ServiceDiscovery {

	var endpoints = []string{"127.0.0.1:2379"}
	ser := NewServiceDiscovery(endpoints)
	ser.InitMyService()

	// 监听的是服务的前置,该服务下的所有服务都会被监听,也会被修改
	_ = ser.WatchService("/buzz")
	go func() {
		for {
			select {
			case <-time.Tick(10 * time.Second):
				log.Println(ser.GetServices())
			}
		}
	}()
	return ser
}

type ServiceRegister struct {
	Cli     *clientv3.Client //etcd client
	LeaseID clientv3.LeaseID //租约ID
	//租约keepalieve相应chan
	KeepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	Key           string //key
	Val           string //value
}

// Close 注销服务
func (s *ServiceRegister) Close() error {
	//撤销租约
	if _, err := s.Cli.Revoke(context.Background(), s.LeaseID); err != nil {
		return err
	}
	fmt.Printf("服务注销成功 key:%s value:%s", s.Key, s.Val)
	return s.Cli.Close()
}
