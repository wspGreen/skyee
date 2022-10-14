package scontext

type HandlerFunc func(session uint32, source uint32, cmd string, params []interface{}) []interface{}

// 多个ctx共用一份proto会覆盖 F Pack
type Protocal struct {
	Name   string
	Type   uint8
	F      HandlerFunc
	UnPack func(rawParams []interface{}) []interface{}          // 收包是解码数据 ， 不同 type ，制定的自己逻辑
	Pack   func(params []interface{}) (rawParams []interface{}) // 发包时编码数据 params:{cmd,...}
}

type IContext interface {
	// DispatchMessage(prototype uint8, m *reflect.Method, paramList []reflect.Value) bool
	GetProtoByType(tp uint8) *Protocal
}
