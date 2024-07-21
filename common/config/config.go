package config

import (
	"common/utils"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type AppConf struct {
	v        *viper.Viper
	filepath string
}

func Read(path string) (*AppConf, error) {
	v := viper.New()
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to read config filepath: [%s], error: %s", path, err)
	}
	return &AppConf{
		v:        v,
		filepath: path,
	}, nil
}

func ReadFromJson(content string) (*AppConf, error) {
	v := viper.New()
	reader := strings.NewReader(content)
	v.SetConfigType("json")
	v.ReadConfig(reader)
	return &AppConf{
		v:        v,
		filepath: "",
	}, nil
}

func NewForUnitTest(paths ...string) *AppConf {
	var path string
	if len(paths) > 0 {
		path = paths[0]
	}

	v := viper.New()
	v.SetConfigFile(path)
	_ = v.ReadInConfig()
	return &AppConf{
		v:        v,
		filepath: path,
	}
}

func (b *AppConf) Set(path string, value interface{}) {
	b.v.Set(path, value)
}

func (b *AppConf) WatchConfig() {
	b.v.WatchConfig()
}

func (b *AppConf) ConfPath() string {
	return b.filepath
}

func (b *AppConf) GetInt(path string, def int) int {
	if b.v.IsSet(path) {
		return b.v.GetInt(path)
	}
	return def
}

func (b *AppConf) GetBool(path string, def bool) bool {
	if b.v.IsSet(path) {
		return b.v.GetBool(path)
	}
	return def
}

func (b *AppConf) GetFloat(path string, def float64) float64 {
	if b.v.IsSet(path) {
		return b.v.GetFloat64(path)
	}
	return def
}

func (b *AppConf) GetStrSlice(path string, def []string) []string {
	if b.v.IsSet(path) {
		return b.v.GetStringSlice(path)
	}
	return def
}

func (b *AppConf) MustGetStrSlice(path string) []string {
	if !b.v.IsSet(path) {
		panic(any(fmt.Sprintf("Cannot get config in %s %s", b.filepath, path)))
	}
	return b.v.GetStringSlice(path)
}

func (b *AppConf) GetString(path string, def string) string {
	if b.v.IsSet(path) {
		return b.v.GetString(path)
	}
	return def
}

func (b *AppConf) GetIntSlice(path string, def []int) []int {
	if b.v.IsSet(path) {
		return b.v.GetIntSlice(path)
	}
	return def
}

func (b *AppConf) Get(path string, def interface{}) interface{} {
	if b.v.IsSet(path) {
		return b.v.Get(path)
	}
	return def
}

func (b *AppConf) MustGetInt(path string) int {
	if !b.v.IsSet(path) {
		panic(any(fmt.Sprintf("Cannot get config in %s %s", b.filepath, path)))
	}
	return b.v.GetInt(path)
}

func (b *AppConf) MustGetInt64(path string) int64 {
	if !b.v.IsSet(path) {
		panic(any(fmt.Sprintf("Cannot get config in %s %s", b.filepath, path)))
	}
	return b.v.GetInt64(path)
}

func (b *AppConf) MustGetFloat(path string) float64 {
	if !b.v.IsSet(path) {
		panic(any(fmt.Sprintf("Cannot get config in %s %s", b.filepath, path)))
	}
	return b.v.GetFloat64(path)
}

func (b *AppConf) MustGetString(path string) string {
	if !b.v.IsSet(path) {
		panic(any(fmt.Sprintf("Cannot get config in %s %s", b.filepath, path)))
	}
	return b.v.GetString(path)
}

func (b *AppConf) GetMap(path string) map[string]interface{} {
	if !b.v.IsSet(path) {
		panic(any(fmt.Sprintf("Cannot get config in %s %s", b.filepath, path)))
	}
	return b.v.GetStringMap(path)
}

func (b *AppConf) MustGetEncryptedString(path string) string {
	cipherString := b.MustGetString(path)
	if cipherString == "" {
		return cipherString
	}
	plainString, err := utils.SimpleDecrypt(cipherString)
	if err != nil {
		panic(any(err))
	}
	return string(plainString)
}
