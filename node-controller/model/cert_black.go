package model

import "time"

const (
	TblUserBlackToken = "tb_user_black_token"
)

type UserBlackToken struct {
	ID         int    `gorm:"column:id;primary_key;auto_increment"`
	UserID     int    `gorm:"uniqueIndex:ix_user_token_user_id_black_token;column:user_id;type:bigint(20);not null"`
	BlackToken string `gorm:"uniqueIndex:ix_user_token_user_id_black_token;column:black_token;type:varchar(255);not null"`
	ExpireAt   time.Time
}

func (UserBlackToken) TableName() string {
	return TblUserBlackToken
}
