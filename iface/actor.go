package iface

type IActor interface {
	Start()
	GetId() uint32

	// 调用函数，参数为 fparams
	// fparams:{...,params:{cmd,...}}
	FunCall(cmd string, fparams ...interface{}) []interface{}
}
