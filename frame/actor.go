package frame

import (
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"

	"github.com/wspGreen/skyee/actor"
	"github.com/wspGreen/skyee/log"
	"github.com/wspGreen/skyee/mpsc"
)

func NewActor(id uint32) *Actor {
	return &Actor{id: id,
		msgs:   mpsc.New(),
		msgCtr: actor.NewMsgCtrl(),
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
}

type IActor interface {
	Init()
}

func (a *Actor) Init(handle IHander) {

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

func (a *Actor) run() {

	for {
		msg := a.msgs.Pop()
		if msg == nil {
			runtime.Gosched()
			continue
		}
		select {
		case <-a.msgCtr.GetMark():
			a.msgCtr.EnableMarkState()
			// msg := a.msgs.Pop()
			// if msg == nil {
			// 	continue
			// }
			Go(func() {
				a.execute(msg.(*Msg))
			})
		}
	}

}

func (a *Actor) Execute(msg *Msg) {
	a.execute(msg)
}

func (a *Actor) execute(msg *Msg) {
	defer func() {
		a.msgCtr.Mark()

		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("Actor Handle cmd[%s] panic: %+v\n%s", msg.CMD, err, debug.Stack()))
		}
	}()

	funcname := msg.CMD

	if !a.hasMethod(funcname) {
		log.Fatalf(" func [%s] has no method ", funcname)
		return
	}

	m, _ := a.refType.MethodByName(funcname)

	paramList := make([]reflect.Value, 0)
	paramList = append(paramList, a.refVal)
	for _, v := range msg.Params {
		paramList = append(paramList, reflect.ValueOf(v))
	}

	m.Func.Call(paramList)

}

func (a *Actor) Send(msg *Msg) {
	a.msgs.Push(msg)

	a.msgCtr.Mark()
}

func (a *Actor) hasMethod(name string) bool {
	_, bEx := a.refType.MethodByName(name)
	return bEx
}
