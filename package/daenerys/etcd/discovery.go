package etcd

import (
	"context"
	"fmt"
	"log"

	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
	clientv3 "go.etcd.io/etcd/clientv3"
	"my.service/go-login/package/config"
)

//ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	cli          *clientv3.Client  //etcd client
	serverList   map[string]string //服务列表
	lock         sync.Mutex
	myService    map[string]struct{} //本服务名称
	serviceMap   map[string]map[string]string
	LoadBanlance map[string]interface{} // 负载均衡算法
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
// map[/buzz/app/go-login/leimingdeAir:127.0.0.1:10006 /buzz/app/go-login/newService:172.0.0.1:10007 /buzz/app/go-login/newService3:172.0.0.1:10008 /buzz/app/go-login/newService4:172.0.0.1:10009 /buzz/app/go-login/newService5:172.0.0.1:10000 /buzz/app/go-login/newService6:172.0.0.1:10001 /buzz/app/go-login/newService7:172.0.0.1:33]
func (s *ServiceDiscovery) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[key] = string(val)

	// 加入到对应的负载均衡器中
	name := GetNameSpaceName(key)
	if name == "" {
		return
	}
	if _, ok := s.LoadBanlance[name]; !ok {
		return
	}

	switch s.LoadBanlance[name].(type) {
	case *RandomBalance:
		s.LoadBanlance[name].(*RandomBalance).Add(val)
	}
	log.Println("put key :", key, "val:", val)
}

//DelServiceList 删除服务地址
func (s *ServiceDiscovery) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// 加入到对应的负载均衡器中
	name := GetNameSpaceName(key)
	if name == "" {
		return
	}
	if _, ok := s.LoadBanlance[name]; !ok {
		return
	}

	switch s.LoadBanlance[name].(type) {
	case *RandomBalance:
		s.LoadBanlance[name].(*RandomBalance).Remove(s.serverList[key])
	}

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

func GetNameSpaceName(serviceName string) string {
	tmp := strings.Split(serviceName, "/")
	if len(tmp) < 4 {
		return ""
	}
	return tmp[1] + "." + tmp[2] + "." + tmp[3]
}

//Close 关闭服务
func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}

func (s *ServiceDiscovery) InitMyService() {
	s.myService = make(map[string]struct{})
	// map[buzz.app.go-login:{} buzz.app.go-user:{}]
	load := make(map[string]interface{})
	for _, v := range config.Conf.ServerClient {
		// 重复配置检查
		if _, ok := s.myService[v.ServiceName]; ok {
			continue
		}
		s.myService[v.ServiceName] = struct{}{}

		// 新的负载均衡配置可以在这里添加
		switch v.Balancetype {
		case Random:
			load[v.ServiceName] = &RandomBalance{}
		}
	}

	s.LoadBanlance = load
}

// 获取当前config所关注的所有服务发现地址
func (s *ServiceDiscovery) GetServicesMap() map[string]map[string]string {
	serviceMap := make(map[string]map[string]string)
	onceMap := make(map[string]string)
	for k, v := range s.serverList {
		tmp := strings.Split(k, "/")
		newServiceName := GetNameSpaceName(k)
		if _, ok := s.myService[newServiceName]; ok {
			onceMap[tmp[4]] = v
			serviceMap[newServiceName] = onceMap
		}
	}
	return serviceMap
}

func InitDiscovery() *ServiceDiscovery {

	var endpoints = []string{"127.0.0.1:2379"}
	ser := NewServiceDiscovery(endpoints)
	ser.InitMyService()

	// 监听的是服务的前置,该服务下的所有服务都会被监听,也会被修改
	_ = ser.WatchService("/buzz")

	ser.serviceMap = ser.GetServicesMap()
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
