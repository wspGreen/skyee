package frame

import (
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/wspGreen/skyee/actor"
	"github.com/wspGreen/skyee/iface"
	"github.com/wspGreen/skyee/scontext"
	"github.com/wspGreen/skyee/slib/mpsc"
	"github.com/wspGreen/skyee/slog"
)

func NewActor(id uint32, ctx scontext.IContext) *Actor {
	return &Actor{id: id,
		msgs:   mpsc.New(),
		msgCtr: actor.NewMsgCtrl(),
		ctx:    ctx,
	}
}

type Actor struct {
	// id int64
	// msgs    chan *net.Packet
	id     uint32
	msgs   *mpsc.Queue
	msgCtr *actor.MsgCtrl
	name   string

	refType reflect.Type
	refVal  reflect.Value
	ctx     scontext.IContext
}

func (a *Actor) Init(handle iface.IHander) {

	a.refType = reflect.TypeOf(handle)
	a.refVal = reflect.ValueOf(handle)
	a.name = GetClassName(a.refType)

	Go(func() {
		a.run()
	})
}

func (a *Actor) Start() {

}

func (a *Actor) Name() string {
	return a.name
}

func (a *Actor) GetId() uint32 {
	return a.id
}

func (a *Actor) run() {
	for {
		// msg := a.msgs.Pop()
		// if msg == nil {
		// 	Sleep(10)
		// 	continue
		// }
		select {
		case <-a.msgCtr.IsEnable():
			a.consume()

		}
	}

}

func (a *Actor) consume() {
	msg := a.msgs.Pop()
	if msg == nil {
		slog.Error("msg state is 1, but msg data is nil !!!")
		return
	}
	Go(func() {
		defer func() {
			a.msgCtr.Disable()
			if !a.msgs.Empty() {
				a.msgCtr.Enable()
			}

			// if err := recover(); err != nil {
			// 	slog.Error(fmt.Sprintf("Actor Handle cmd[%s] panic: %+v\n%s", msg.CMD, err, debug.Stack()))
			// }
		}()
		a.execute(msg.(*SkyeeMsg))
	})
}

func (a *Actor) Execute(msg *SkyeeMsg) {
	a.execute(msg)
}

func (a *Actor) execute(msg *SkyeeMsg) {
	defer func() {

		if err := recover(); err != nil {
			cmd, _ := msg.Params[0].(string)
			slog.Error(fmt.Sprintf("Actor:%s Handle cmd[%s] panic: %+v\n%s", a.name, cmd, err, debug.Stack()))
			// slog.Error(fmt.Sprintf("Actor Handle cmd[%s] \n", msg.CMD))
		}
	}()

	// funcname := msg.CMD

	// if !a.hasMethod(funcname) {
	// 	slog.Warn(" func [%s] has no method ", funcname)
	// 	return
	// }

	a.dispatchMessage(msg.ProtoType, msg)

}

func (a *Actor) dispatchMessage(prototype uint8, msg *SkyeeMsg) {
	proto := a.ctx.GetProtoByType(prototype)
	if proto == nil {
		return
	}

	if prototype == PTYPE_RESPONSE {

		return
	}

	// funcname := msg.CMD

	// m, _ := a.refType.MethodByName(funcname)

	params := proto.UnPack(msg.Params)

	// paramList := make([]reflect.Value, 0)
	// paramList = append(paramList, a.refVal)
	// for _, v := range params {
	// 	paramList = append(paramList, reflect.ValueOf(v))
	// }
	cmd := params[0].(string)
	params = params[1:] // remove funcname form params

	r := proto.F(msg.Session, msg.Source, cmd, params)
	if r != nil {
		command, ok := r[0].(string)
		if ok {
			a.suspend(command)
		}
	}
}

func (a *Actor) suspend(command string, params ...interface{}) {
	if command == "RETURN" {
		// session:=0
		// source := 0
		// frame.PTYPE_RESPONSE
		// Send(source, "RESPONSE", "", params...)
	}
}

func (a *Actor) Send(msg *SkyeeMsg) {
	a.msgs.Push(msg)

	a.msgCtr.Enable()
}

func (a *Actor) hasMethod(name string) bool {
	_, bEx := a.refType.MethodByName(name)
	return bEx
}

// fparams:{...,params:{...}}
func (a *Actor) FunCall(cmd string, fparams ...interface{}) []interface{} {
	l := len(fparams)
	params, ok := fparams[l-1].([]interface{})
	if !ok {
		slog.Error(" fparams[%s] not  type : []interface{}", l)
		return nil
	}
	fparams = fparams[:l-1] // remove params form fparams
	// funcname := params[0].(string)
	// params = params[1:] // remove funcname form params

	m, ok := a.refType.MethodByName(cmd)
	if !ok {
		slog.Warn(" func [%s] has no method ", cmd)
		return nil
	}

	paramList := []reflect.Value{a.refVal}
	for _, v := range fparams {
		paramList = append(paramList, reflect.ValueOf(v))
	}

	for i := 0; i < len(params); i++ {
		paramList = append(paramList, reflect.ValueOf(params[i]))
	}

	// for _, v := range params {
	// 	paramList = append(paramList, reflect.ValueOf(v))
	// }

	ret := m.Func.Call(paramList)

	if ret != nil {

	}
	return nil
}
