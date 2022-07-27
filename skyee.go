package skyee

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/wspGreen/skyee/component"
	"github.com/wspGreen/skyee/frame"
	"github.com/wspGreen/skyee/log"
)

var stopChanForSys = make(chan os.Signal, 1)

var components = &component.Components{}
var contextId uint32

/////////////////////////////////////////

type OptionFunc func(c *frame.SkyeeContext, params *OptionParam)
type OptionParam struct {
	Id uint32
}

func NewService(h frame.IHander, options ...OptionFunc) uint32 {
	id := atomic.AddUint32(&contextId, 1)
	ctx := frame.NewContext(id, h)

	ac := frame.NewActor(id)
	ac.Init(h)
	// ac.Start()

	ctx.SetActor(ac)

	frame.AddContext(ctx)

	for _, e := range options {
		o := new(OptionParam)
		e(ctx, o)
	}

	start()

	log.Println(fmt.Sprintf("start service %s[%d]", ac.Name(), id))

	return ctx.GetId()
}

func start() {
	components.Start()
}

func UniqueService() {

}

func SetWebSocket(ctxId uint32) {
	// actorid 关联到 socket 才能收 网络数据

	AddComponent(frame.NewNetServerComp("ws", ":9001", ctxId))
}

func SetHttp(ctxId uint32, addr string) {
	// actorid 关联到 socket 才能收 网络数据

	AddComponent(frame.NewNetServerComp("http", addr, ctxId))
}

func AddComponent(c component.IComponent) {
	components.AddComponent(c)
}

func Send(id uint32, typename string, cmd string, params ...interface{}) {
	ctx := frame.GetContext(id)
	msg := frame.NewMsg(typename, cmd, params)

	ctx.GetActor().Send(msg)
}

func WaitForSystemExit() {

	signal.Notify(stopChanForSys, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-stopChanForSys:
		Stop()
	}

	// fmt.Printf("server【%s】 exit ------- signal:[%v]", args[1], s)
	// fmt.Printf("server exit ------- signal:[%v]", s)
}

func Stop() {
	log.Println("Server Stop")
	close(stopChanForSys)
}
