package base

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func init() {
	Log.Info("开始连接数据库...")
	createConnect()
}

func createConnect() {
	var err error
	Db, err = gorm.Open("mysql", Conf.Dbconfig)
	if err != nil {
		Log.Error("数据库连接失败，10秒钟后重试...")
		time.Sleep(10 * time.Second)
		createConnect()
	} else {
		Log.Info("数据库连接成功 ^_^")
		Db.DB().SetConnMaxLifetime(100 * time.Second) //最大连接周期，超过时间的连接就close
		Db.DB().SetMaxOpenConns(100)                  //设置最大连接数
		Db.DB().SetMaxIdleConns(0)                    //设置闲置连接数
		Db.SingularTable(true)                        //设置全局表名禁用复数
		//指定表前缀，修改默认表名
		gorm.DefaultTableNameHandler = func(b *gorm.DB, defaultTableName string) string {
			return Conf.DbPrefix + defaultTableName
		}
	}
}
