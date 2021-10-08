package controller

import (
	"encoding/json"
	"fmt"
	"manage/base"
	"manage/model"
	"time"
)

type user_group_permission struct {
	model.UserGroupPermission
}

var UserGroupPermission user_group_permission

/*
根据用户组ID和roterid检查是否有权限
*/
func (p *user_group_permission) PermissionCheckByGroupId(group_id int, roter_id int) bool {
	var d_p model.UserGroupPermission
	base.Db.Where("group_id=? AND menu_id=?", group_id, roter_id).Find(&d_p)
	if d_p.Id > 0 {
		return true
	} else {
		return false
	}
}

type usergroup struct {
	model.UserGroup
}

var UserGroup usergroup

/*
获取分组列表，全部，提供给一些select使用
*/
func (ug *usergroup) FullList() *[]model.UserGroup {
	var _list []model.UserGroup
	base.Db.Find(&_list)
	return &_list
}

/*
用户组分页列表
*/
func (ug *usergroup) List(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}

	var total int
	base.Db.Model(&model.UserGroup{}).Count(&total)

	p := &model.Page{}
	currentPage, pageSize, startNum := p.Handle(cmd, total)

	list := &[]model.UserGroup{}
	base.Db.Table(fmt.Sprintf("%suser_group ug", base.Conf.DbPrefix)).Select("ug.*").Limit(pageSize).Offset(startNum).Find(&list)

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
用户分组信息修改
*/
func (ug *usergroup) Modify(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	ugjson, _ := json.Marshal(cmd.Data)
	var _usergroup model.UserGroup
	err := json.Unmarshal(ugjson, &_usergroup)
	if err != nil {
		res.Code = 50000
		res.Message = "参数错误"
	} else {
		if _usergroup.GroupName == "" {
			res.Code = 50000
			res.Message = "用户分组名称必须输入"
		} else {
			if _usergroup.Id > 0 {
				var _ug model.UserGroup
				// 系统分组的属性锁定
				base.Db.Where("id=?", _usergroup.Id).First(&_ug)
				if _ug.System {
					_usergroup.System = true
				}
				_usergroup.UpdateTime = int(time.Now().Unix())
				base.Db.Save(&_usergroup)
			} else {
				_usergroup.CreateTime = int(time.Now().Unix())
				_usergroup.UpdateTime = int(time.Now().Unix())
				base.Db.Create(&_usergroup)
			}
			res.Code = 20000
			res.Message = ""
		}
	}
	return res
}

/*
分组删除

1.删除分组数据
2.将关联的用户 的 分组id设置成0
3.删除分组权限绑定 中的数据
*/
func (ug *usergroup) Remove(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	type delQuery struct {
		Id int `json:"id"`
	}

	_del_json, _ := json.Marshal(cmd.Data)
	var _del_query delQuery
	err := json.Unmarshal(_del_json, &_del_query)
	if err != nil {
		res.Code = 50000
		res.Message = "参数错误"
	} else {
		// 系统分组不能删除
		var _g model.UserGroup
		base.Db.Where("id=?", _del_query.Id).First(&_g)
		if _g.System {
			res.Code = 50000
			res.Message = "系统用户组，不能删除"
		} else {
			// 如果删除的分组是当前用户分组，不允许删除
			if _del_query.Id == cmd.State.GroupId {
				res.Code = 50000
				res.Message = "当前登录用户的分组不能删除"
			} else {
				// 设置分组的用户gropu_id 为 0
				var _u model.User
				base.Db.Model(&_u).Where("group_id=?", _del_query.Id).Update("group_id", 0)
				// 删除分组的绑定数据
				var _up model.UserGroupPermission
				base.Db.Where("group_id=?", _del_query.Id).Delete(&_up)
				// 删除分组数据
				var _ug model.UserGroup
				base.Db.Where("id=?", _del_query.Id).Delete(&_ug)

				res.Code = 20000
				res.Message = ""
			}
		}

	}
	return res
}

/*
获取权限配置列表以及指定分组的权限列表

1.获取菜单列表，生成权限树
2.根据group_id获取分组的权限列表
3.保存分组的权限数据
*/
func (ug *usergroup) Permission(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	type queryObj struct {
		Action   string `json:"action"`
		GroupId  int    `json:"group_id"`
		AuthList []int  `json:"auth_list"`
	}
	queryjson, _ := json.Marshal(cmd.Data)
	var _query queryObj
	err := json.Unmarshal(queryjson, &_query)
	if err != nil {
		res.Code = 50000
		res.Message = "参数错误"
	} else {
		if _query.Action == "save" {
			menu_auth_save_for_group(_query.GroupId, _query.AuthList)
		} else {
			data := make(map[string]interface{})
			data["auth_tree"] = menu_tree_list(0)
			data["auth_list"] = menu_auth_for_group(_query.GroupId)
			res.Data = data
		}
	}
	return res
}

type menuTree struct {
	Id       int        `json:"id"`
	Label    string     `json:"label"`
	Children []menuTree `json:"children"`
}

// 菜单列表生成elemui 的tree数据源
func menu_tree_list(parentid int) *[]menuTree {
	list := []model.Menu{}
	base.Db.Table(fmt.Sprintf("%smenu m", base.Conf.DbPrefix)).Where(fmt.Sprintf("m.parent_id=%v", parentid)).Select("m.*").Order("m.sort asc").Find(&list)
	menulist := []menuTree{}
	for i := 0; i < len(list); i++ {
		list_item := list[i]
		menuItem := &menuTree{
			Id:    list_item.Id,
			Label: list_item.Name,
		}
		pchildren := menu_tree_list(list_item.Id)
		menuItem.Children = *pchildren
		menulist = append(menulist, *menuItem)
	}
	return &menulist
}

// 获取指定分组已经有的权限列表
func menu_auth_for_group(groupid int) []int {
	list := []model.UserGroupPermission{}
	base.Db.Table(fmt.Sprintf("%suser_group_permission ugp", base.Conf.DbPrefix)).Where(fmt.Sprintf("ugp.group_id=%v", groupid)).Select("ugp.*").Find(&list)

	var auth_list []int
	for i := 0; i < len(list); i++ {
		menu_item := list[i]
		auth_list = append(auth_list, menu_item.MenuId)
	}
	return auth_list
}

/*
根据用户组保存权限树信息
*/
func menu_auth_save_for_group(groupid int, list []int) {
	////保存权限
	for i := 0; i < len(list); i++ {
		var _c int
		base.Db.Table(fmt.Sprintf("%suser_group_permission ugp", base.Conf.DbPrefix)).Where("group_id=? and menu_id=?", groupid, list[i]).Count(&_c)
		if _c == 0 {
			base.Db.Create(&model.UserGroupPermission{GroupId: groupid, MenuId: list[i]})
		}
	}
	// 删除被取消的权限，如果menuid不再list中的删除
	var _l model.UserGroupPermission
	base.Db.Where("group_id=? and  menu_id not in (?)", groupid, list).Delete(&_l)
}
