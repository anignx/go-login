package dao

type UserAuths struct {
	ID           uint `gorm:"primarykey"`
	UserID       uint `gorm:"unique"`
	IdentityType int
	Identifier   string
	Credential   string
	CreateTime   int64 `gorm:"autoCreateTime"`
	UpdateTime   int64 `gorm:"autoUpdateTime"`
}

func (UserAuths) TableName() string {
	return "user_auths"
}
