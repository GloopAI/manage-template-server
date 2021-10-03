package model

// 子模块返回实体
type CommandResult struct {
	Code    int
	Message string
	Data    interface{} //map[string]interface{}
	// DateTest interface{}
}
