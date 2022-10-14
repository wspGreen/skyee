package frame

import (
	"sync"

	"github.com/wspGreen/skyee/iface"
	"github.com/wspGreen/skyee/scontext"
	"github.com/wspGreen/skyee/slog"
)

var skyeeContexts = map[interface{}]*SkyeeContext{}
var contextSync sync.Mutex

func AddContext(ctx *SkyeeContext) {
	contextSync.Lock()
	skyeeContexts[ctx.GetId()] = ctx
	contextSync.Unlock()
}

func AddContextByName(ctx *SkyeeContext) {
	contextSync.Lock()
	skyeeContexts[ctx.GetName()] = ctx
	contextSync.Unlock()
}

func GetContext(id interface{}) *SkyeeContext {
	return skyeeContexts[id]
}

/////////////////////////////////////////
// Contxt
// 这个类更像是一个VM？
type SkyeeContext struct {
	id     uint32
	handle iface.IHander // handle logic
	ac     *Actor

	protos     map[uint8]*scontext.Protocal
	protosName map[string]*scontext.Protocal
}

func NewContext(id uint32, handle iface.IHander) *SkyeeContext {
	return &SkyeeContext{
		id:         id,
		handle:     handle,
		protos:     make(map[uint8]*scontext.Protocal),
		protosName: make(map[string]*scontext.Protocal),
	}
}

func (ctx *SkyeeContext) GetId() uint32 {
	return ctx.id
}

func (ctx *SkyeeContext) GetName() string {
	return ctx.ac.Name()
}

func (ctx *SkyeeContext) SetActor(ac *Actor) {
	ctx.ac = ac
}

func (ctx *SkyeeContext) GetActor() *Actor {
	return ctx.ac
}

func (ctx *SkyeeContext) RegisterProtocol(proto *scontext.Protocal) bool {

	ctx.protos[proto.Type] = proto
	if proto.Name == "" {
		slog.Fatal("RegisterProtocol error ")
		return false
	}
	ctx.protosName[proto.Name] = proto
	return true
}

func (ctx *SkyeeContext) GetProto(name string) *scontext.Protocal {
	return ctx.protosName[name]
}

func (ctx *SkyeeContext) GetProtoByType(tp uint8) *scontext.Protocal {
	return ctx.protos[tp]
}

// func (ctx *SkyeeContext) DispatchMessage(prototype uint8, m *reflect.Method, paramList []reflect.Value) bool {
// 	proto := ctx.protos[prototype]
// 	if proto == nil {
// 		return false
// 	}

// 	r := proto.F(m, paramList)
// 	return r
// }

// 设置处理函数
func (ctx *SkyeeContext) Dispatch(protoname string, fun scontext.HandlerFunc) {
	proto := ctx.protosName[protoname]
	if proto == nil {
		return
	}

	proto.F = fun
}

type SkyeeMsg struct {
	ProtoType uint8
	// ProtoName string
	CMD     string
	Session uint32
	Source  uint32
	Params  []interface{} // 编码后的参数
}

func NewSkyeeMsg(prototype uint8, Session uint32, Source uint32, params []interface{}) *SkyeeMsg {
	return &SkyeeMsg{
		ProtoType: prototype,
		Session:   Session,
		Source:    Source,
		// CMD:     cmd,
		Params: params,
	}
}
