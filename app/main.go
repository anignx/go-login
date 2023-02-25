package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"my.service/go-login/conf"
	logging "my.service/go-login/package/BackendPlatform/logging"
	myconfig "my.service/go-login/package/config"
	"my.service/go-login/package/daenerys"
	"my.service/go-login/server/http"
	"my.service/go-login/service"
)

func init() {
	configS := flag.String("config", "config/config.toml", "Configuration file")
	flag.Parse()

	myconfig.Init(
		*configS,
	)
	daenerys.Init(daenerys.ConfigPath(*configS))
}

func main() {

	logging.LogConf()

	cfg, err := conf.Init()
	if err != nil {
		logging.Logger.Fatalf("service config init error %s", err)
	}

	service := service.New(cfg)

	http.Init(service, cfg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-sigChan
		log.Printf("get a signal %s\n", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Println("gmu.social.chat_service server exit now...")
			return
		case syscall.SIGHUP:
		default:
		}
	}
}
