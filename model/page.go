package model

import (
	"manage/base"

	"github.com/goinggo/mapstructure"
)

type Page struct {
	Page     int `json:"page"`
	Pagesize int `json:"pagesize"`
}

func (p *Page) Handle(cmd *Command, total int) (CurrentPage int, PageSize int, StartNum int) {
	var _p Page
	_p.Page = 1
	_p.Pagesize = base.Conf.PageSize
	err := mapstructure.Decode(cmd.Data, &_p)
	if err != nil {
		_p.Page = 1
		_p.Pagesize = base.Conf.PageSize
	}

	startNum := (_p.Page - 1) * _p.Pagesize
	return _p.Page, _p.Pagesize, startNum
}
