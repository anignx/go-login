package proxy

import (
	"gorm.io/gorm"
	"my.service/go-login/package/daenerys"
)

type SQL struct {
	name []string
}

// Client继承了*gorm.DB的所有方法, 详细的使用方法请参考:
// http://gorm.io/docs/connecting_to_the_database.html
type Client struct {
	*gorm.DB
}

func InitSQL(name ...string) *SQL {
	if len(name) == 0 {
		return nil
	}
	return &SQL{name}
}

func (s *SQL) Master(name ...string) *gorm.DB {
	var gName string
	if len(name) == 0 {
		gName = s.name[0]
	} else {
		gName = name[0]
	}
	return daenerys.Defult.SQLClient(gName)
}

func (s *SQL) Slave(name ...string) *gorm.DB {
	var gName string
	if len(name) == 0 {
		gName = s.name[0]
	} else {
		gName = name[0]
	}
	return daenerys.Defult.SQLClient(gName)
}
