package main

import (
	"manage/base"
	"manage/roter"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 设置跨域
	cf := cors.DefaultConfig()
	cf.AllowMethods = []string{"*"}
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

	r.Run() // listen and serve on 0.0.0.0:8080
	base.Log.Info("服务启动，端口:8080")
}
