package controller

import (
	"encoding/json"
	"fmt"
	"manage/base"
	"manage/model"

	"github.com/go-playground/validator/v10"
	"github.com/goinggo/mapstructure"
)

type menu struct {
	model.Menu
}

var Menu *menu

/*
根据command 获取 Roter信息
*/
func (r *menu) InfoByCommand(cmd *model.Command) *model.Menu {
	res := &model.Menu{Id: 0, Name: "", RoterCommand: ""}
	base.Db.Where("roter_command = ?", cmd.Command).First(&res)
	return res
}

/*
获取用户菜单
*/
func (r *menu) AuthMenu(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: make(map[string]interface{})}
	res.Data = menu_list_by_parent_id_with_userid(0, cmd.State.GroupId)
	return res
}

// 根据parent_menu_id 获取下级的菜单列表
func menu_list_by_parent_id_with_userid(parentid int, user_group_id int) *[]model.Menu {
	list := []model.Menu{}
	base.Db.Table(fmt.Sprintf("%smenu m", base.Conf.DbPrefix)).Joins(fmt.Sprintf("left join %suser_group_permission ugp on m.id = ugp.menu_id", base.Conf.DbPrefix)).Where(fmt.Sprintf("ugp.group_id=%v and m.parent_id=%v", user_group_id, parentid)).Select("m.*").Order("m.sort asc").Find(&list)
	menulist := []model.Menu{}
	for i := 0; i < len(list); i++ {
		list_item := list[i]
		menuItem := &model.Menu{
			Id:        list_item.Id,
			ParentId:  list_item.ParentId,
			Name:      list_item.Name,
			Icon:      list_item.Icon,
			Component: list_item.Component,
			Hidden:    list_item.Hidden,
		}
		pchildren := menu_list_by_parent_id_with_userid(list_item.Id, user_group_id)
		menuItem.Children = *pchildren
		menulist = append(menulist, *menuItem)
	}
	return &menulist
}

/*
获取菜单列表
{"page":1,"pagesize":10}
*/
func (m *menu) List(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	type queryObj struct {
		ParentId int `json:"parentid"`
	}

	query := queryObj{ParentId: 0}
	err := mapstructure.Decode(cmd.Data, &query)
	if err != nil {
		fmt.Println(err)
		query.ParentId = 0
	}
	var total int
	base.Db.Model(&model.Menu{}).Where("parent_id=?", query.ParentId).Count(&total)

	p := &model.Page{}
	currentPage, pageSize, startNum := p.Handle(cmd, total)

	list := &[]model.Menu{}
	base.Db.Table(fmt.Sprintf("%smenu m", base.Conf.DbPrefix)).Select("m.*").Where("parent_id=?", query.ParentId).Order("m.sort asc").Limit(pageSize).Offset(startNum).Find(&list)

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
保存菜单信息
*/
func (m *menu) Modify(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	menu_json, _ := json.Marshal(cmd.Data)
	// fmt.Println(string(menu_json))
	var r_menu model.Menu
	err := json.Unmarshal([]byte(menu_json), &r_menu)
	if err != nil {
		res.Code = 50000
		res.Message = err.Error()
	} else {
		validate := validator.New()
		err := validate.Struct(r_menu)
		if err != nil {
			errobj, _ := err.(validator.ValidationErrors)
			res.Code = 50000
			res.Message = errobj[0].Error()
		} else {
			var _m model.Menu
			var _c int
			// var _c 、= 0
			// 检查名称是否存在
			base.Db.Where("name=? and id!=?", r_menu.Name, r_menu.Id).Find(&_m).Count((&_c))
			if _c > 0 {
				res.Code = 50000
				res.Message = "菜单名称不能重复"
			} else {
				// 检查roter_command
				base.Db.Where("roter_command=? and id!=?", r_menu.RoterCommand, r_menu.Id).Find(&_m).Count(&_c)
				if _c > 0 {
					res.Code = 50000
					res.Message = "RoterCommand 不能重复"
				} else {
					// 检查 Component
					base.Db.Where("component=? and id !=?", r_menu.Component, r_menu.Id).Find(&_m).Count(&_c)
					if _c > 0 {
						res.Code = 50000
						res.Message = "Component 不能重复"
					} else {
						// 检查通过之后，做入库
						// id 大于 0 更新
						if r_menu.Id > 0 {
							// 系统菜单的属性锁定
							base.Db.Where("id=?", r_menu.Id).First(&_m)
							if _m.System {
								r_menu.System = true
							}
							base.Db.Save(&r_menu)
						} else {
							base.Db.Create(&r_menu)
						}
						res.Code = 20000
						res.Message = ""
					}
				}
			}

		}
	}

	return res
}

/*
删除菜单信息

1.删除菜单记录 _menu
2.删除菜单权限绑定 _user_group_permission
3.检查有没有子菜单，如果有的话，不能删除
*/
func (m *menu) Remove(cmd *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	type delQuery struct {
		Id int `json:"id"`
	}
	queryJson, _ := json.Marshal(cmd.Data)
	var _del_query delQuery
	err := json.Unmarshal(queryJson, &_del_query)
	if err != nil {
		res.Code = 50000
		res.Message = "参数错误"
	} else {
		// 有系统菜单标识的不能删除
		var _del_menu model.Menu
		base.Db.Where("id=?", _del_query.Id).First(&_del_menu)
		if _del_menu.System {
			res.Code = 50000
			res.Message = "系统菜单不能删除"
		} else {
			// 检查是否有子菜单
			var _m model.Menu
			var _c int
			base.Db.Where("parent_id=?", _del_query.Id).Find(&_m).Count(&_c)
			if _c > 0 {
				res.Code = 50000
				res.Message = "包含子菜单，不能删除"
			} else {
				// 删除菜单表
				base.Db.Where("id=?", _del_query.Id).Delete(&_m)
				// 删除权限绑定表
				var _mp model.UserGroupPermission
				base.Db.Where("menu_id=?", _del_query.Id).Delete(&_mp)
				res.Code = 20000
				res.Message = ""
			}
		}

	}

	return res
}
