package main

import (
	"fmt"
	"io"
	"manage/base"
	"manage/roter"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
}

func main() {
	// gin 的debug设置
	if !base.Conf.IsDebug {
		logfile, err := os.Create("gin_http.log")
		if err != nil {
			fmt.Println("Could not create log file")
		}
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.MultiWriter(logfile)
	}

	r := gin.Default()
	// 设置跨域
	cf := cors.DefaultConfig()
	cf.AllowMethods = []string{base.Conf.AllowMethods}
	cf.AddAllowHeaders("X-Token")
	cf.AllowAllOrigins = true
	cf.AllowCredentials = true
	cf.MaxAge = 24 * time.Hour
	r.Use(cors.New(cf))

	r.POST("/manage/api", func(c *gin.Context) {
		roter.Manage.Receive(c)
	})

	r.POST("/api", func(c *gin.Context) {
		roter.Api.Recive(c)
	})

	r.Run(base.Conf.BindPort) // listen and serve on 0.0.0.0:8080
	base.Log.Info("服务启动，端口:8080")
}
