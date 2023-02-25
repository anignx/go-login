package dao

import (
	"my.service/go-login/package/BackendPlatform/logging"
)

type User struct {
	ID         uint `gorm:"primarykey"`
	Nickname   string
	Avatar     string
	CreateTime int64 `gorm:"autoCreateTime"`
	UpdateTime int64 `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "user"
}

func (d *Dao) GetUserByID(id uint) (*User, error) {
	logging.Logger.Infof("get user by id: %d", id)
	user := &User{}
	err := d.db.Master().Table("user").Where("id = ?", id).First(user).Error
	return user, err
}
