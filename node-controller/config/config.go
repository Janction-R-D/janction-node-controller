package config

import (
	"common/config"
	"common/db"
	"node-controller/common/cache"
	"node-controller/common/cache/mem"
	"node-controller/common/certificate"
	"node-controller/dao/query"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

func MustInitConfig(conf *config.AppConf) {
	initApp(conf)
	InitMemCache()
	initMysqlDb(conf)
	initJwt()
	initENV(conf)
}

var (
	AppHost          string
	AppPort          int
	DB               *db.DB
	CasbinModelPath  string
	CasbinPolicyPath string
	ENV              string
	Endpoint         string
	TokenAddress     string
	TotalAmount      int
)
var (
	PaymentAddress  string
	StartTime       int64
	EndTime         int64
	OriginPrice     float64
	DiscountPrice   float64
	FundraisingGoal float64
)

var Jwt *certificate.JwtCertProvider
var CacheInstance cache.Cache

func initApp(conf *config.AppConf) {
	AppHost = conf.MustGetString("app.host")
	AppPort = conf.MustGetInt("app.port")
}

func InitMemCache() {
	CacheInstance = mem.New(5*time.Minute, 30*time.Minute)
}

func initMysqlDb(conf *config.AppConf) {
	DB, _ = db.Connect(conf, logger.Default.LogMode(logger.Info), 10)
	query.SetDefault(DB.DB)
}

func initJwt() {
	Jwt = certificate.NewJwtCertProvider(certificate.JwtKey("uih2ia", "cdjwiu234u29h"),
		certificate.JwtCertProviderWithLoginDur(24*time.Hour),
	)
}

func parseIntervalAndDuration(conf *config.AppConf, path, defaultInterval string, defaultDuration time.Duration) (string, time.Duration) {
	interval := conf.GetString(path, defaultInterval)
	duration, err := time.ParseDuration(interval)
	if err != nil {
		logrus.WithError(err).WithField(path, interval).Error("filed to parse duration")
		return defaultInterval, defaultDuration
	}
	return interval, duration
}

func initENV(conf *config.AppConf) {
	ENV = conf.GetString("env", "online")
}
