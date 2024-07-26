package postgres

import (
	"fmt"
	"janction/setting"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(cfg *setting.PostgresConfig) (err error) {
	var dsn string

	if cfg.User == "" && cfg.Password == "" {
		dsn = fmt.Sprintf(
			"host=%s port=%d dbname=%s sslmode=disable search_path=%s",
			cfg.Host,
			cfg.Port,
			cfg.DB,
			cfg.Schema,
		)
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable search_path=%s",
			cfg.Host,
			cfg.Port,
			cfg.DB,
			cfg.User,
			cfg.Password,
			cfg.Schema,
		)
	}

	if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		return err
	}

	if err = db.AutoMigrate(

	); err != nil {
		return err
	}

	return nil
}

func Close() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
