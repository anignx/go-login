package dao

type UserStatus struct {
	ID         uint `gorm:"primarykey"`
	UserID     uint `gorm:"unique"`
	Status     int
	CreateTime int64 `gorm:"autoCreateTime"`
	UpdateTime int64 `gorm:"autoUpdateTime"`
}

func (UserStatus) TableName() string {
	return "user_status"
}
