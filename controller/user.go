package controller

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"manage/base"
	"manage/model"
	"time"

	"github.com/goinggo/mapstructure"
	"github.com/google/uuid"
)

type user struct {
	model.User
}

var User *user

func init() {

}

/*
用户登录

{"username":"xxx","password":""}
*/
func (u *user) Login(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	var r_user model.User
	err := mapstructure.Decode(cmd.Data, &r_user)
	if err != nil {
		res.Code = 50000
		res.Message = "账号或者密码错误"
	} else {
		var d_user model.User
		base.Db.Where("username = ?", r_user.Username).First(&d_user)
		if d_user.Id == 0 {
			res.Code = 50000
			res.Message = "账号或者密码错误"
		} else {
			d_pwd := d_user.Password
			r_pwd_byte := []byte(r_user.Password)
			var r_pwd string
			r_pwd = fmt.Sprintf("%x", md5.Sum(r_pwd_byte))
			// fmt.Println(r_pwd)
			if d_pwd != r_pwd {
				res.Code = 50000
				res.Message = "账号或者密码错误"
			} else {
				_token := fmt.Sprintf("%v", uuid.New().String())
				d_user.Token = _token
				base.Db.Model(&d_user).Where("id=?", d_user.Id).Update("token", _token)

				data := make(map[string]interface{})
				data["token"] = _token
				res.Code = 20000
				res.Message = ""
				res.Data = data
			}
		}
	}
	return res
}

/*
根据token获取用户信息

{"token":"xxx-xxx-xxx-xxx"}
*/
func (u *user) InfoByToken(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	var r_user model.User
	err := mapstructure.Decode(cmd.Data, &r_user)
	if err != nil {
		res.Code = 50008
		res.Message = "登录超时 Login failed, unable to get user details."
	} else {
		var d_user model.User
		base.Db.Where("token=?", r_user.Token).First(&d_user)
		if d_user.Id == 0 {
			res.Code = 50008
			res.Message = "Login failed, unable to get user details."
		} else {
			d_user.Password = ""
			d_user.Token = ""
			res.Data = d_user
		}
	}
	return res
}

/*
退出登录
*/
func (u *user) LoginOut(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "success", Data: nil}
	return res
}

/*
 获取用户分页列表

 {"page":1,"pagesize":10}
*/
func (u *user) List(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: make(map[string]interface{})}

	var total int
	base.Db.Model(&model.User{}).Count(&total)

	p := &model.Page{}
	currentPage, pageSize, startNum := p.Handle(cmd, total)

	list := &[]model.UserExt{}
	base.Db.Table(fmt.Sprintf("%suser u", base.Conf.DbPrefix)).Joins(fmt.Sprintf("left join %suser_group ug on ug.id = u.group_id", base.Conf.DbPrefix)).Select("u.*,ug.group_name").Limit(pageSize).Offset(startNum).Find(&list)

	pageData := make(map[string]interface{})
	pageData["total"] = total
	pageData["page"] = currentPage
	pageData["pagesize"] = pageSize
	pageData["list"] = list
	items := make(map[string]interface{})
	items["items"] = pageData
	res.Data = items
	return res
}

/*
修改用户信息，username不能修改，token不能修改
{"id":1,"name":"evan","password":"xxxxx"}
*/
func (u *user) Modify(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	type queryObj struct {
		Action string     `json:"action"`
		Data   model.User `json:"data"`
	}
	queryjson, _ := json.Marshal(cmd.Data)
	var query queryObj
	err := json.Unmarshal(queryjson, &query)
	if err != nil {
		res.Code = 50000
		res.Message = err.Error()
	} else {
		if query.Action == "save" {
			m := modify_save(&query.Data)
			res.Code = m.Code
			res.Message = m.Message
		} else {
			result := make(map[string]interface{})
			result["group_list"] = UserGroup.FullList()
			res.Data = result
		}
	}
	return res
}

/*
用户自己修改密码
*/
func (u *user) PasswordReset(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	type reset struct {
		Old        string `json:"old"`
		New        string `json:"new"`
		ConfirmNew string `json:"confirm_new"`
	}
	reset_json, _ := json.Marshal(cmd.Data)
	var query reset
	err := json.Unmarshal(reset_json, &query)
	if err != nil {
		res.Code = 50000
		res.Message = "参数错误"
		return res
	}
	if query.Old == "" {
		res.Code = 50000
		res.Message = "需要输入旧密码"
		return res
	}

	var _u model.User
	base.Db.Model(&_u).Where("id=?", cmd.State.Id).First(&_u)
	if _u.Id == 0 {
		res.Code = 50000
		res.Message = "用户参数错误"
		return res
	}
	_old_pass_byte := []byte(query.Old)
	_old_pass := fmt.Sprintf("%x", md5.Sum(_old_pass_byte))
	if _old_pass != _u.Password {
		res.Code = 50000
		res.Message = "旧密码错误"
		return res
	}

	if query.New == "" {
		res.Code = 50000
		res.Message = "需要输入新密码"
		return res
	}
	if query.New != query.ConfirmNew {
		res.Code = 50000
		res.Message = "两次新密码输入需要一致"
		return res
	}

	// 开始修改密码
	_new_passwd_byte := []byte(query.New)
	_new_passwd := fmt.Sprintf("%x", md5.Sum(_new_passwd_byte))
	__u := model.User{
		Password:   _new_passwd,
		UpdateTime: int(time.Now().Unix()),
	}
	base.Db.Model(&__u).Where("id=?", cmd.State.Id).Update(&__u)
	return res
}

/*
编辑用户信息保存
*/
func modify_save(user *model.User) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	// fmt.Println(user)
	if user.Username == "" {
		res.Code = 50000
		res.Message = "UserName 必须输入"
	} else if user.NickName == "" {
		res.Code = 50000
		res.Message = "NickName 必须输入"
	} else if user.GroupId == 0 {
		res.Code = 50000
		res.Message = "Group 必须选择"
	} else {

		if user.Id == 0 {
			_n_pass_byte := []byte("123456")
			if user.Password != "" {
				_n_pass_byte = []byte(user.Password)
			}
			user.Password = fmt.Sprintf("%x", md5.Sum(_n_pass_byte))
			u := model.User{
				Username:   user.Username,
				NickName:   user.NickName,
				Password:   user.Password,
				GroupId:    user.GroupId,
				CreateTime: int(time.Now().Unix()),
				UpdateTime: int(time.Now().Unix()),
			}
			base.Db.Create(&u)
		} else {
			u := model.User{
				Username:   user.Username,
				NickName:   user.NickName,
				GroupId:    user.GroupId,
				UpdateTime: int(time.Now().Unix()),
			}
			if user.Password != "" {
				_n_pass_byte := []byte(user.Password)
				u.Password = fmt.Sprintf("%x", md5.Sum(_n_pass_byte))
			}
			base.Db.Model(&u).Where("id=?", user.Id).Update(&u)
		}

	}
	return res
}

/*
删除用户
*/
func (u *user) Remove(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	queryjson, _ := json.Marshal(cmd.Data)
	var query model.User
	err := json.Unmarshal(queryjson, &query)
	if err != nil {
		res.Code = 50000
		res.Message = "参数错误"
	} else {
		var _u model.User
		base.Db.Where("id=?", query.Id).First(&_u)
		if _u.System {
			res.Code = 50000
			res.Message = "系统账户不能删除"
		} else {
			// 当前登录账户不能删除
			if query.Id == cmd.State.Id {
				res.Code = 50000
				res.Message = "当前登录账户不能删除"
			} else {
				base.Db.Where("id=?", query.Id).Delete(&query)
			}
		}
	}
	return res
}
