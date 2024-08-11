package model

import "time"

type JWT struct {
	Token     string    `json:"token" gorm:"type:varchar(300)"`
	ExpiredAt time.Time `json:"expired_at"`
}
