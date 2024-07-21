package db

import (
	"common/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLConf struct {
	Driver      string
	User        string
	Password    string
	Host        string
	Port        uint16
	DBName      string
	MaxOpenConn int
	MaxIdleConn int
	MaxLifeTm   int
}

type ConnOption struct {
	Interval time.Duration
	MaxTimes int
	Logging  func(err error)
}

type OptionModels func(db *gorm.DB) (*gorm.DB, error)

type DB struct {
	*gorm.DB
}

func NewDB(config *MySQLConf, l logger.Interface) (*DB, error) {
	conn, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local&parseTime=true", config.User, config.Password, config.Host, config.Port, config.DBName),
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 l,
	})
	if err != nil {
		return nil, err
	}
	db, err := conn.DB()
	if config.MaxOpenConn > 0 {
		db.SetMaxOpenConns(config.MaxOpenConn)
	}
	if config.MaxIdleConn > 0 {
		db.SetMaxIdleConns(config.MaxIdleConn)
	}
	if config.MaxLifeTm > 0 {
		db.SetConnMaxLifetime(time.Duration(config.MaxLifeTm) * time.Second)
	}
	return &DB{DB: conn}, nil
}

func Retry(executor func() (*DB, error), opt ConnOption) (*DB, error) {
	var err error
	var db *DB
	for i := 0; i < opt.MaxTimes; i++ {
		db, err = executor()
		if err == nil {
			return db, nil
		} else {
			if opt.Logging != nil {
				opt.Logging(err)
			}
		}
		time.Sleep(opt.Interval)
	}
	return nil, err
}

func Connect(conf *config.AppConf, l logger.Interface, retryTimes int) (*DB, error) {
	db, err := Retry(func() (*DB, error) {
		dbConfig := MySQLConf{
			Host:        conf.MustGetString("db.host"),
			Port:        uint16(conf.GetInt("db.port", 3306)),
			DBName:      conf.MustGetString("db.dbname"),
			User:        conf.MustGetString("db.user"),
			Password:    conf.MustGetString("db.password"),
			MaxOpenConn: conf.GetInt("db.max_conn", 10),
			MaxIdleConn: conf.GetInt("db.max_idle_conn", 2),
			MaxLifeTm:   conf.GetInt("db.max_life_tm", 300),
		}
		db, err := NewDB(&dbConfig, l)
		if err != nil {
			return nil, err
		}
		return db, nil
	}, ConnOption{
		Interval: time.Second,
		MaxTimes: retryTimes,
		Logging: func(err error) {
			log.Println("connect db db fail: ", err)
		},
	})
	if err != nil {
		panic(any("connect to database failed, abort"))
	}
	return db, nil
}

// OptionPaging 分页
func (m *DB) OptionPaging(page, pageSize int) OptionModels {
	return func(gormDB *gorm.DB) (*gorm.DB, error) {
		if page <= 0 || pageSize <= 0 {
			return nil, fmt.Errorf("page or page size < 0")
		}
		return gormDB.Limit(pageSize).Offset((page - 1) * pageSize), nil
	}
}

// OptionLike 模糊搜索
func (m *DB) OptionLike(field, pattern string) OptionModels {
	return func(gormDB *gorm.DB) (*gorm.DB, error) {
		key := fmt.Sprintf("%s LIKE ?", field)
		value := fmt.Sprintf("%%%s%%", pattern)
		return gormDB.Where(key, value), nil
	}
}

// OptionOrder 排序
func (m *DB) OptionOrder(ordering string) OptionModels {
	return func(gormDB *gorm.DB) (*gorm.DB, error) {
		var order, field string
		if ordering[0] == '-' {
			order = "desc"
			field = ordering[1:]
		} else if ordering[0] == '+' {
			order = "asc"
			field = ordering[1:]
		} else {
			order = "asc"
			field = ordering[1:]
		}

		search := fmt.Sprintf("%s %s", field, order)

		return gormDB.Order(search), nil
	}
}

// Transaction
func (m *DB) Transaction(things ...func(tx2 *gorm.DB) error) error {
	return m.DB.Transaction(func(tx *gorm.DB) error {
		for _, t := range things {
			if err := t(tx); err != nil {
				return err
			}
		}
		return nil
	})
}

type Page struct {
	Offset int
	Limit  int
}
