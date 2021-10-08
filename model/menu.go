package model

type Menu struct {
	Id           int    `gorm:"id" json:"id"`
	Name         string `gorm:"name" json:"name,omitempty" validate:"required"`
	Icon         string `gorm:"icon" json:"icon"  validate:"required"`
	RoterCommand string `gorm:"roter_command" json:"roter_command,omitempty"  validate:"required"`
	Component    string `gorm:"component" json:"component"  validate:"required"`
	ParentId     int    `gorm:"parent_id" json:"parent_id"`
	Hidden       bool   `gorm:"hidden" json:"hidden"`
	System       bool   `gorm:"system" json:"system"`
	Sort         int    `gorm:"sort" json:"sort" `
	Note         string `gorm:"note" json:"note"`
	Children     []Menu `json:"children,omitempty"`
	CreateTime   int    `gorm:"create_time" json:"create_time"`
	UpdateTime   int    `gorm:"update_time" json:"update_time"`
}
