package roter

import (
	"encoding/json"
	"fmt"
	"manage/base"
	"manage/controller"
	"manage/model"

	"github.com/gin-gonic/gin"
)

func init() {
}

type manage struct {
}

var Manage *manage

// 接收post
func (r *manage) Receive(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read((buf))
	var cmd model.Command
	err := json.Unmarshal(buf[0:n], &cmd)

	if err != nil {
		if base.Conf.IsDebug {
			base.Log.Error(fmt.Sprintf("post data解析失败 %s", string(buf[0:n])))
		}
		c.JSON(200, gin.H{"code": 50000, "message": "post data is error", "data": ""})
	} else {
		base.Log.Info(fmt.Sprintf("接收到参数 command:%s data:%s", cmd.Command, cmd.Data))
		cmd.State.Token = c.Request.Header.Get("x-token")

		if Manage.permissionCheck(&cmd) {
			res := Manage.commandExc(&cmd)
			c.JSON(200, gin.H{"code": res.Code, "message": res.Message, "data": res.Data})
		} else {
			c.JSON(200, gin.H{"code": 50000, "message": "没有权限", "data": make(map[string]interface{})})
		}
	}
}

// 命令处理路由
func (r *manage) commandExc(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	switch cmd.Command {
	case "user.login":
		res = controller.User.Login(cmd)
	case "user.login.out":
		res = controller.User.LoginOut(cmd)
	case "user.info.by.token":
		res = controller.User.InfoByToken(cmd)
	case "sys.menu.auth":
		res = controller.Menu.AuthMenu(cmd)
	case "sys.menu.list":
		res = controller.Menu.List(cmd)
	case "sys.menu.modify":
		res = controller.Menu.Modify(cmd)
	case "sys.menu.remove":
		res = controller.Menu.Remove(cmd)
	case "user.list":
		res = controller.User.List(cmd)
	case "user.modify":
		res = controller.User.Modify(cmd)
	case "user.remove":
		res = controller.User.Remove(cmd)
	case "user.password.reset":
		res = controller.User.PasswordReset(cmd)
	case "user.group.list":
		res = controller.UserGroup.List(cmd)
	case "user.group.modify":
		res = controller.UserGroup.Modify(cmd)
	case "user.group.remove":
		res = controller.UserGroup.Remove(cmd)
	case "user.group.permission":
		res = controller.UserGroup.Permission(cmd)
	case "test":
		res = controller.App.Test(cmd)
	default:
		res.Code = 50000
		res.Message = "command not found!!"
	}
	if res.Data == nil {
		res.Data = make(map[string]interface{})
	}
	return res
}

// 权限检查
func (r *manage) permissionCheck(cmd *model.Command) bool {
	if cmd.Command == "user.login" || cmd.Command == "user.info.by.token" || cmd.Command == "user.login.out" {
		return true
	} else {
		c_cmd := &model.Command{}
		c_data := make(map[string]interface{})
		c_data["token"] = cmd.State.Token
		c_cmd.Data = c_data
		_u := controller.User.InfoByToken(c_cmd).Data
		// fmt.Println(_u)
		// 检查用户是否登录
		if _u != nil {
			user := _u.(model.User)
			// 用户已登录
			if user.Id != 0 {
				cmd.State.IsLogin = true
				cmd.State.Id = user.Id
				cmd.State.NickName = user.NickName
				cmd.State.GroupId = user.GroupId
				// 免权限验证路由
				if cmd.Command == "sys.menu.auth" || cmd.Command == "user.password.reset" {
					return true
				} else {
					d_roter := controller.Menu.InfoByCommand(cmd)
					// 菜单模块权限检查
					if d_roter.Id > 0 {
						if controller.UserGroupPermission.PermissionCheckByGroupId(user.GroupId, d_roter.Id) {
							return true
						} else {
							return false
						}
					} else {
						return false
					}
				}
			} else {
				return false
			}
		} else {
			return false
		}
	}
}
