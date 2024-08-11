package postgres

import (
	"janction/model"
	"time"

	"gorm.io/gorm"
)

func SaveOrReplaceJWT(token string, expiredAt time.Time) error {
	var jwt model.JWT
	result := db.First(&jwt)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			newJWT := model.JWT{
				Token:     token,
				ExpiredAt: expiredAt,
			}
			return db.Create(&newJWT).Error
		}
		return result.Error
	}

	jwt.Token = token
	jwt.ExpiredAt = expiredAt
	return db.Save(&jwt).Error
}

func GetJWT() (*model.JWT, error) {
	var jwt model.JWT
	result := db.First(&jwt)
	if result.Error != nil {
		return nil, result.Error
	}
	return &jwt, nil
}
