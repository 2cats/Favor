package models

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/gcfg.v1"
)

func panicWhenError(err error) {
	if err != nil {
		panic(err)
	}
}

var ConfigFileName string

const (
	DefaultConfFile = "server.conf"
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
		UploadDir         string
	}
}
var Config = &_Config.Common

func ReadConfiguration() {
	err := gcfg.ReadFileInto(&_Config, ConfigFileName)
	panicWhenError(err)
	if _, err := os.Stat(Config.UploadDir); os.IsNotExist(err) {
		fmt.Printf("UploadDir Not Exist,Creating...")
		err = os.Mkdir(Config.UploadDir, 0700)
		panicWhenError(err)
		fmt.Printf("%s Created", Config.UploadDir)
	}
	fmt.Printf("Config:%+v\n", *Config)
}
