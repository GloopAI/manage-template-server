package model

type ApiCommand struct {
	Method    string
	Appkey    string
	Timestamp string
	Sign      string
	Data      map[string]interface{}
}
