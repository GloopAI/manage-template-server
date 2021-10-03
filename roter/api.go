package roter

import (
	"encoding/json"
	"fmt"
	"manage/base"
	"manage/model"

	"github.com/gin-gonic/gin"
)

type api struct {
	model.ApiCommand
}

var Api *api

func (a *api) Recive(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read((buf))
	var cmd model.ApiCommand
	err := json.Unmarshal(buf[0:n], &cmd)
	if err != nil {
		if base.Conf.IsDebug {
			base.Log.Error(fmt.Sprintf("post data解析失败 %s", string(buf[0:n])))
		}
		c.JSON(200, gin.H{"code": 50000, "message": "post data is error", "data": ""})
	} else {
		base.Log.Info(fmt.Sprintf("接收到参数 method:%s appkey:%s data:%s", cmd.Method, cmd.Appkey, cmd.Data))
		res := Api.commandExc(&cmd)
		c.JSON(200, gin.H{"code": res.Code, "message": res.Message, "data": res.Data})
	}
}

func (a *api) commandExc(cmd *model.ApiCommand) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	if cmd.Appkey == "" {
		res.Code = 50000
		res.Message = "缺少参数appkey"
	} else if cmd.Sign == "" {
		res.Code = 50000
		res.Message = "缺少参数sign"
	} else if cmd.Timestamp == "" {
		res.Code = 50000
		res.Message = "缺少参数timestamp"
	} else {
		switch cmd.Method {

		default:
			res.Code = 50000
			res.Message = "method not found!!"
		}
	}

	if res.Data == nil {
		res.Data = make(map[string]interface{})
	}
	return res
}
