package skyee

import (
	"math"
	"reflect"
	"sync/atomic"

	"github.com/wspGreen/skyee/component"
	"github.com/wspGreen/skyee/frame"
	"github.com/wspGreen/skyee/iface"
	"github.com/wspGreen/skyee/scontext"
	"github.com/wspGreen/skyee/slog"
)

var components = component.NewComponents()
var contextId uint32

/////////////////////////////////////////

type OptionFunc func(c *frame.SkyeeContext, params *OptionParam)
type OptionParam struct {
	Id uint32
}

func NewService(h iface.IHander, options ...OptionFunc) uint32 {

	return newService(h, false, options...)
}

func newService(h iface.IHander, isUnique bool, options ...OptionFunc) uint32 {
	ctx := newContext(h, isUnique)
	for _, e := range options {
		o := new(OptionParam)
		e(ctx, o)
	}
	// 组件是给处理类用的先准备完成
	start()

	// 初始化处理类
	h.Init(ctx.GetActor())

	// slog.Info("start service %s[%d] unique:%v", ctx.GetActor().Name(), ctx.GetId(), isUnique)
	slog.Info("start service %s[%d] unique:%v addr:%p", ctx.GetActor().Name(), ctx.GetId(), isUnique, h)

	return ctx.GetId()
}

func UniqueService(h iface.IHander, options ...OptionFunc) uint32 {

	name := frame.GetClassName(reflect.TypeOf(h))
	ctx := frame.GetContext(name)
	if ctx != nil {
		return ctx.GetId()
	}
	return newService(h, true, options...)
}

func initProto(ctx *frame.SkyeeContext) {
	ctx.RegisterProtocol(&scontext.Protocal{
		Name:   "cmd",
		Type:   frame.PTYPE_CMD,
		UnPack: UnPack,
		Pack:   Pack,
	})

}

func Start(f func()) {
	frame.Go(f)
	waitForSystemExit()
}

func start() {

	components.Start()
}

func SetFileLog(path string) slog.ISLog {
	l := slog.Log()
	l.SetLogger(&frame.FileLogger{Path: path, Ln: true, MaxSize: frame.LOG_FILE_MAX_SIZE})
	return l
}

func newContext(h iface.IHander, isUnique bool) *frame.SkyeeContext {

	if atomic.LoadUint32(&contextId) == math.MaxUint32 {
		atomic.StoreUint32(&contextId, 0)
	}
	id := atomic.AddUint32(&contextId, 1)
	ctx := frame.NewContext(id, h)

	ac := frame.NewActor(id, ctx)
	ac.Init(h)
	// ac.Start()

	ctx.SetActor(ac)

	frame.AddContext(ctx)
	if isUnique {
		frame.AddContextByName(ctx)
	}

	initProto(ctx)

	return ctx
}

func SetWebSocket(ctxId uint32, addr string) {
	// actorid 关联到 socket 才能收 网络数据

	AddComponent(frame.NewNetServerComp("ws", addr, ctxId))
}

func SetHttp(ctxId uint32, addr string) {
	// actorid 关联到 socket 才能收 网络数据

	AddComponent(frame.NewNetServerComp("http", addr, ctxId))
}

func SetNats(h iface.IBrokerHandler) {
	AddComponent(frame.NewNatsComp(h))
}

func Nats() *frame.NatsComp {
	return GetComponent("NatsComp").(*frame.NatsComp)
}

func Net() *frame.NetServerComp {
	return GetComponent("NetServerComp").(*frame.NetServerComp)
}

func AddComponent(c iface.IComponent) {
	components.AddComponent(c)
}

func GetComponent(name string) iface.IComponent {
	return components.GetComponent(name)
}

// 发送给actor或集群
// params : {cmd,...}|{source,cmd}
func Send(id interface{}, protoname string, params ...interface{}) {

	if isCluster(id) {
		// sendCluster(id, cmd, params...)
	} else {

		send(id, protoname, params...)

	}

}

func Call(id interface{}, protoname string, cmd string, params ...interface{}) {
	// send(id, protoname, cmd, params...)
}

func send(id interface{}, protoname string, params ...interface{}) uint32 {
	ctx := frame.GetContext(id)
	if ctx == nil {
		slog.Error("context id:%v is not exist", id)
		return 0
	}
	proto := ctx.GetProto(protoname)
	if proto == nil {
		slog.Error("protoname is not exist : %s", protoname)
		return 0
	}
	if proto.F == nil {
		slog.Error("protoname:%v is not handle function", proto.Name)
		return 0
	}
	session := uint32(0)
	source, ok := params[0].(uint32)
	if ok {
		params = params[1:]
	}
	cmd, ok := params[0].(string)
	if !ok {
		slog.Info("cmd is error :%v", cmd)
		return 0
	}
	msg := frame.NewSkyeeMsg(proto.Type, session, source, proto.Pack(params))

	ctx.GetActor().Send(msg)
	return session
}

func Ret(ret []interface{}) []interface{} {
	if ret == nil {
		return nil
	}

	return []interface{}{"RETURN"}
}

/*
	发送到集群
	node 集群点 id/类型
	actorid 集群点的服务
*/
func SendCluster(node string, actorid uint32, cmd string, params ...interface{}) {
	// Nats().Send(node,)
}

func isCluster(id interface{}) bool {
	switch id.(type) {
	case uint32:
		return false
	case string:
		// slog.Error("send not support string arg!! ")
		return false
	default:
		slog.Error("send arg type error !! ")
		return false
	}
}

func RegisterProtocol(ctxid uint32, proto *scontext.Protocal) {
	ctx := frame.GetContext(ctxid)
	ctx.RegisterProtocol(proto)
}

func Dispatch(ctxid uint32, protoname string, fun scontext.HandlerFunc) {
	ctx := frame.GetContext(ctxid)
	ctx.Dispatch(protoname, fun)
}

// 解码
func UnPack(rawParams []interface{}) []interface{} {
	return rawParams
}

// 编码
func Pack(params []interface{}) (pParams []interface{}) {
	pParams = params
	return
}

func waitForSystemExit() {

	frame.WaitForSystemExit()

	// fmt.Printf("server【%s】 exit ------- signal:[%v]", args[1], s)
	// fmt.Printf("server exit ------- signal:[%v]", s)
}
