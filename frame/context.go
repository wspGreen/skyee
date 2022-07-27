package frame

import "sync"

var skyeeContexts = map[uint32]*SkyeeContext{}
var contextSync sync.Mutex

func AddContext(ctx *SkyeeContext) {
	contextSync.Lock()
	skyeeContexts[ctx.GetId()] = ctx
	contextSync.Unlock()
}

func GetContext(id uint32) *SkyeeContext {
	return skyeeContexts[id]
}

/////////////////////////////////////////
// Contxt
// 这个类更像是一个VM？
func NewContext(id uint32, handle IHander) *SkyeeContext {
	return &SkyeeContext{
		id:     id,
		handle: handle,
	}
}

type SkyeeContext struct {
	id     uint32
	handle IHander
	ac     *Actor
}

func (ctx *SkyeeContext) GetId() uint32 {
	return ctx.id
}

func (ctx *SkyeeContext) SetActor(ac *Actor) {
	ctx.ac = ac
}

func (ctx *SkyeeContext) GetActor() *Actor {
	return ctx.ac
}

type Msg struct {
	Typeid   uint8
	Typename string
	CMD      string
	Params   []interface{}
}

func NewMsg(typename string, cmd string, params []interface{}) *Msg {
	return &Msg{
		Typeid:   0,
		Typename: typename,
		CMD:      cmd,
		Params:   params,
	}
}
