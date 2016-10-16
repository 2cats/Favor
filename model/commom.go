package models

import (
	"strconv"
	"time"
)

const (
	UPLOAD_DIR = "public/upload/"
	AdminUser  = "admin"
)

func GetFilePrefix() string {
	return strconv.FormatInt(time.Now().Unix(), 10) + "_"
}
