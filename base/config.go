package base

import (
	"encoding/json"
	"os"
)

// 配置文件结构体
type configuration struct {
	Dbconfig string
	DbPrefix string
	IsDebug  bool
	PageSize int
}

var Conf configuration

func init() {
	// log.SetFormatter(&log.JSONFormatter{})
	// log.SetLevel(log.WarnLevel)

	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&Conf)
	if err != nil {
		Log.Error("配置文件读取失败")
	} else {
		Log.Info("配置文件读取成功")
	}
}
