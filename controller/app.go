package controller

import (
	"manage/model"
)

type app struct {
}

var App app

func (a *app) Test(c *model.Command) *model.CommandResult {
	res := &model.CommandResult{Code: 20000, Message: "", Data: nil}
	// fmt.Println(c.State.Name)

	return res
}
