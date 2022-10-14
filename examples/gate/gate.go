package gate

import (
	"fmt"
	"strconv"

	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/frame"
	"github.com/wspGreen/skyee/iface"
	"github.com/wspGreen/skyee/slib"
	"github.com/wspGreen/skyee/slib/consts"
	"github.com/wspGreen/skyee/slog"
)

// var Gated frame.Service = NewGate

// var Gated = NewGate()

type Gate struct {
	MsgServer
	id uint32

	Username_map map[string]*User
}

type User struct {
	agent uint32
}

func NewGate() *Gate {
	return &Gate{
		Username_map: make(map[string]*User),
		MsgServer: MsgServer{
			handshake:   map[uint32]string{},
			user_online: map[string]*ConnectInfo{},
			connection:  map[uint32]*ConnectInfo{},
		},
	}
}

func (g *Gate) Init(a iface.IActor) {
	g.MsgServer.Init(a, g)

	skyee.Dispatch(
		a.GetId(),
		"cmd",
		func(session uint32, source uint32, cmd string, params []interface{}) []interface{} {
			ret := a.FunCall(cmd, params)

			return skyee.Ret(ret)
		},
	)

	skyee.Nats().Subscribe(strconv.Itoa(g.getSvrId()))
	g.id = a.GetId()
	// slog.Info(" gate init id:%v", g.id)

}

func (g *Gate) Start() {
	// slog.Info("gate Start addr:%p", g)
}

// call by login server
func (g *Gate) LoginHandler(uid string, account string) {
	// uid 账号id
	// account 账号名
	username := username(uid, g.getSvrId())
	g.Username_map[username] = &User{agent: 1}
	g.MsgServer.login(username, account)
	slog.Info("user login ok uid:%s,acc:%s", uid, account)
}

func username(uid string, svrid int) string {
	return fmt.Sprintf("%s-%d", uid, svrid)
}

func (g *Gate) getSvrId() int {
	return consts.SERVER_TYPE_GATE
}

func (g *Gate) requestHandler(username string, s *frame.Session, p *frame.Packet) {
	cmd := ""
	server_type := getProtoRouterServerType(cmd)
	if server_type != consts.SERVER_TYPE_GATE {
		srv_id := getPlayerBindServer(username)
		// 转发给其他服务器
		user := g.Username_map[username]
		slog.Debug("Gate To Server (playerid:%d,serverid:%d)", username, srv_id)
		source := g.id
		slib.ForwardServerMessage(strconv.Itoa(srv_id), user.agent, source, "OnServerMessage", username, p.Data)
	} else {
		onClientRequest(cmd)
	}
}

// 接收服务间信息
/*
	两种情况收到：
	1. 其他服务发给gate的包，在gate处理 (包来源 nats )
	2. 其他服务让gate转发client的包，只需转发给client （包来源client）
*/
func (g *Gate) OnServerMessage(pid string, data []byte) {
	// slog.Info(" send OnServerMessage ok! ")

	// 如果是路由过来的Response包的话  直接发给客户端
	cmd := ""
	if GET_MESSAGE_TYPE(cmd) == consts.Msg_Type_PlayerResponse {
		g.sendResponseClient(pid, data)
		return

	}
	fun := getServerCmd(cmd)
	fun()

}

func (g *Gate) sendResponseClient(pid string, data []byte) {
	g.MsgServer.respone(pid, data)

}

func GET_MESSAGE_TYPE(cmd string) int {
	return consts.Msg_Type_PlayerResponse
}

// 让给gate处理网络数据
func onClientRequest(cmd string) {
	// fun := getClientCmd(cmd)
	// fun()
}

func getProtoRouterServerType(cmd string) int {
	return consts.SERVER_TYPE_GAME
}

func getServerCmd(cmd string) func() {
	// panic("unimplemented")
	return nil
}

func getClientCmd(cmd string) func() {
	// panic("unimplemented")
	return nil
}

func (g *Gate) Move(a int) {
	slog.Info(" send Move %d", a)
}

func (g *Gate) Move1(a uint32) {
	slog.Info(" send Move1 %d", a)
}

func (g *Gate) Attck(name string) {
	slog.Info(" send Attck %s", name)
}

func (g *Gate) NewUniqueSvr(uid uint32) {
	slog.Info(" [id:%v] NewUniqueSvr %v addr:%p", g.id, uid, g)
}
