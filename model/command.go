package model

// 接收post的json实体
type Command struct {
	Command string      `json:"command"`
	Data    interface{} `json:"data"`
	Token   string
	State   UserExt
}
