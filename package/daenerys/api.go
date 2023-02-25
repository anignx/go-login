package daenerys

import (
	"fmt"
	"strings"

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
			logging.Logger.Errorf("sql init err, conf: %v", conf)
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
			logging.Logger.Errorf("sql init err, name: %s, err: %s", sqlConf.Name, err.Error())
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
