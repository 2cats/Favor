package models

import (
	"fmt"
	"strconv"
	"time"

	"gopkg.in/gcfg.v1"
)

func panicWhenError(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	UPLOAD_DIR     = "public/upload/"
	ConfigFileName = "server.conf"
	AdminUser      = "admin"
)

func GetFilePrefix() string {
	return strconv.FormatInt(time.Now().Unix(), 10) + "_"
}

var _Config struct {
	Common struct {
		Port              string
		Listen            string
		MaxUploadFileSize int
		DBFile            string
	}
}
var Config = &_Config.Common

func ReadConfiguration() {
	err := gcfg.ReadFileInto(&_Config, ConfigFileName)
	panicWhenError(err)
	fmt.Printf("Config:%+v\n", *Config)
}
