package dao

import (
	"my.service/go-login/package/BackendPlatform/logging"
)

type User struct {
	ID         uint   `gorm:"primarykey"`
	Iphone     string `gorm:"unique"`
	Nickname   string
	Avatar     string
	CreateTime int64 `gorm:"autoCreateTime"`
	UpdateTime int64 `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "user"
}

func (d *Dao) GetUserByID(userId uint) (*User, error) {
	logging.Logger.Infof("get user by id: %d", userId)
	user := &User{}
	err := d.db.Master().Table("user").Where("id = ?", userId).First(user).Error
	return user, err
}

// CreateUser 创建用户
func (d *Dao) CreateUser(user *User, identityType int, credential string) (int, error) {
	logging.Logger.Infof("create user: %v", user)
	var userId int
	tx := d.db.Master().Begin()
	defer func() {
		if err := recover(); err != nil {
			logging.Logger.Errorf("create user error: %v", err)
			d.db.Master().Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return userId, err
	}
	if err := tx.Table("user").Create(user).Error; err != nil {
		tx.Rollback()
		return userId, err
	}
	logging.Logger.Infof("create user success: new ID:%d", user.ID)

	userAuths := &UserAuths{
		UserID:       user.ID,
		IdentityType: identityType,
		Identifier:   user.Iphone,
		Credential:   credential,
	}
	if err := tx.Table("user_auths").Create(userAuths).Error; err != nil {
		tx.Rollback()
		return userId, err
	}
	logging.Logger.Infof("create user_auths success: new ID:%d", userAuths.ID)

	UserStatus := &UserStatus{
		UserID: user.ID,
		Status: 0,
	}

	if err := tx.Table("user_status").Create(UserStatus).Error; err != nil {
		tx.Rollback()
		return userId, err
	}
	logging.Logger.Infof("create user_status success: new ID:%d", UserStatus.ID)
	return int(user.ID), tx.Commit().Error
}
