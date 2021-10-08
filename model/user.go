package model

type User struct {
	Id         int    `gorm:"id" json:"id"`
	Username   string `gorm:"username" json:"username,omitempty"`
	Password   string `gorm:"password" json:"password,omitempty"`
	Token      string `gorm:"token" json:"token,omitempty"`
	NickName   string `gorm:"nick_name" json:"nick_name"`
	GroupId    int    `gorm:"group_id" json:"group_id"`
	System     bool   `gorm:"system" json:"system"`
	CreateTime int    `gorm:"create_time" json:"create_time"`
	UpdateTime int    `gorm:"update_time" json:"update_time"`
}

type UserExt struct {
	User
	GroupName string `json:"group_name,omitempty"`
	IsLogin   bool   `json:"is_login,omitempty"`
}
