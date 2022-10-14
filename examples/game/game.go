package game

import (
	"fmt"
	"strconv"

	"github.com/bitly/go-simplejson"
	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/iface"
	"github.com/wspGreen/skyee/slib"
	"github.com/wspGreen/skyee/slib/consts"
	"github.com/wspGreen/skyee/slog"
)

var Gamed = NewGame()

type Game struct {
	// slib.BaseServer
	source uint32
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Init(a iface.IActor) {
	skyee.Dispatch(
		a.GetId(),
		"cmd",
		func(session uint32, source uint32, cmd string, params []interface{}) []interface{} {
			a.FunCall(cmd, source, params)

			return nil
		},
	)
	skyee.Nats().Subscribe(strconv.Itoa(g.GetSvrId()))
}

func (g *Game) GetSvrId() int {
	return consts.SERVER_TYPE_GAME
}

func (g *Game) OnServerMessage(source uint32, pid string, data []byte) {

	js, err := simplejson.NewJson(data)
	if err != nil {
		slog.Error(err)
	}

	proto, _ := js.Get("proto").String()
	slog.Info(proto)
	g.source = source
	///////////////////////////////////////
	// game
	// cmd := ""
	if getClientCmd(proto) {
		g.onClientRequest(proto, pid, data)
		return
	}

	fun := getServerCmd(proto)
	fun()
}

// 处理客户端发来的包
func (g *Game) onClientRequest(proto string, pid string, data []byte) {

	if proto == "createrole" {
		Handle_CreateRole(g.source, pid, data)
	}
	// 给对应agent处理？，在返回数据给gate
	// Handle_XXX(pid, data)
}

func getServerCmd(cmd string) func() {
	panic("unimplemented")
}

// 是否为客户端的包
func getClientCmd(cmd string) bool {
	return true
}

func Handle_CreateRole(source uint32, pid string, data []byte) {
	js, err := simplejson.NewJson(data)
	if err != nil {
		slog.Error(err)
	}
	name, _ := js.Get("rolename").String()
	s := fmt.Sprintf("{rolename:%s}", name)
	resq := []byte(s)
	slib.SendResponseByGate(pid, source, resq)
}
