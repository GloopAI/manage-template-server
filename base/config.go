package base

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

// 配置文件结构体
type configuration struct {
	Dbconfig     string
	DbPrefix     string
	IsDebug      bool
	PageSize     int
	AllowMethods string
	BindPort     string
}

var Conf configuration

func init() {
	log.SetFormatter(&log.TextFormatter{})

	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&Conf)
	if err != nil {
		Log.Error("配置文件读取失败")
	} else {
		Log.Info("配置文件读取成功")
		//
		if Conf.IsDebug {
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetLevel(log.ErrorLevel)
		}

	}
}
