package setting

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Port      int    `mapstructure:"port"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`

	UrlConfig      *UrlConfig      `mapstructure:"url"`
	AuthConfig     *AuthConfig     `mapstructure:"auth"`
	LogConfig      *LogConfig      `mapstructure:"log"`
	PostgresConfig *PostgresConfig `mapstructure:"postgres"`
	RedisConfig    *RedisConfig    `mapstructure:"redis"`
}

type UrlConfig struct {
	PointSystem string `mapstructure:"point_system"`
}

type AuthConfig struct {
	JwtExpire int `mapstructure:"jwt_expire"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type PostgresConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Schema       string `mapstructure:"schema"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

var Config = new(AppConfig)

func Init() (err error) {
	if len(os.Args) < 2 {
		log.Fatalln("Please provide the path to the config file as a command line argument.")
		return
	}

	configFilePath := os.Args[1]

	viper.SetConfigFile(configFilePath)
	if err = viper.ReadInConfig(); err != nil {
		log.Panicf("Fatal error config file: %s\n", err)
	}

	if err = viper.Unmarshal(Config); err != nil {
		log.Fatalln("Viper unmarshal failed:", err)
	}

	return
}
