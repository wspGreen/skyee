package gate

import (
	"fmt"

	"github.com/bitly/go-simplejson"
	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/frame"
	"github.com/wspGreen/skyee/iface"
	"github.com/wspGreen/skyee/scontext"
	"github.com/wspGreen/skyee/slib/consts"
	"github.com/wspGreen/skyee/slog"
)

// 放接受网络数据的接口

// type UserAuth struct {
// 	username string
// }

type ConnectInfo struct {
	sessionid uint32
	ip        string
	username  string
}

type MsgServer struct {
	// slib.BaseServer
	handshake map[uint32]string

	// 已连接的玩家信息，包含未在gate验证的玩家（即 只是通过login验证）
	user_online map[string]*ConnectInfo

	// 已连接的玩家信息，只包含通过login和gate验证的玩家
	// 这里里面的连接发的信息才能被成功发给服务处理
	connection map[uint32]*ConnectInfo

	// request_handler func(s *frame.Session, p *frame.Packet)
	gate *Gate
}

func (msg *MsgServer) Init(a iface.IActor, g *Gate) {
	msg.gate = g
	skyee.RegisterProtocol(a.GetId(), &scontext.Protocal{
		Name: "socket",
		Type: frame.PTYPE_SOCKET,
		Pack: func(params []interface{}) (rawParams []interface{}) {
			cmd := params[0].(string)
			s := params[1].(*frame.Session)
			slog.Info("Pack session %d", s.Id())
			if len(params) > 2 {
				data := params[2].([]byte)
				p := &frame.Packet{Data: data}
				rawParams = []interface{}{cmd, s, p}
			} else {
				rawParams = []interface{}{cmd, s}
			}
			return
		},
		UnPack: func(rawparams []interface{}) []interface{} {
			// s := rawparams[0].(*frame.Session)
			// p := rawparams[1].(*frame.Packet)
			// slog.Info("Unpack %d %v", s.Id(), string(p.Data))
			return rawparams
		},

		F: func(session uint32, source uint32, cmd string, params []interface{}) []interface{} {
			ret := a.FunCall(cmd, params)
			if ret != nil {

			}
			return nil
		},
	})

}

func (msg *MsgServer) login(username string, secret string) {
	msg.user_online[username] = &ConnectInfo{
		username: username,
	}
}

func (msg *MsgServer) OnClientRequest(s *frame.Session, p *frame.Packet) {

	slog.Info("OnClientRequest:%s", string(p.Data))
	js, err := simplejson.NewJson(p.Data)
	if err != nil {
		slog.Fatal("解析数据失败 %v", err)
		return
	}

	// LoginHandler 要由login服调用，现在模拟登录成功先放这边
	t, _ := js.Get("type").String()
	if t == "login" {
		uid, _ := js.Get("uid").String()
		account, _ := js.Get("account").String()
		msg.gate.LoginHandler(uid, account)

		responeCode(s, 0)
		return
	}

	// gate验证
	addr, ok := msg.handshake[s.Id()]
	if ok {
		msg.auth(s, p.Data, addr)
		delete(msg.handshake, s.Id())
		responeCode(s, 0)
	} else {
		msg.request(s, p)
	}

	// 不验证
	// msg.request(s, p)
}

func responeCode(s *frame.Session, code int) {
	s.Write([]byte(fmt.Sprintf("{code:%d}", code)))
}

func (msg *MsgServer) request(s *frame.Session, p *frame.Packet) {
	c, ok := msg.connection[s.Id()]
	if !ok {
		// 后面处理
		responeCode(s, 1)
		return
	}

	msg.gate.requestHandler(c.username, s, p)

}

func (msg *MsgServer) respone(pid string, data []byte) {
	cinfo, ok := msg.user_online[pid]
	if !ok {
		slog.Error("user:%s is null", pid)
		return
	}
	s := skyee.Net().GetSession(cinfo.sessionid)
	if s == nil {
		slog.Error("session id is nil , %d", pid)
		return
	}
	s.Write(data)
}

// 第一次连接到gate需要验证一次
func (msg *MsgServer) auth(s *frame.Session, data []byte, addr string) {
	js, err := simplejson.NewJson(data)
	if err != nil {
		slog.Fatal("解析数据失败 %v", err)
		return
	}
	uid, _ := js.Get("uid").String()
	username := username(uid, msg.gate.getSvrId())

	u, ok := msg.user_online[username]
	if !ok {
		return
	}
	u.sessionid = s.Id()
	msg.connection[s.Id()] = u
	// 验证成功后
	// auth_ok_handler()
}

func (msg *MsgServer) OnOpen(s *frame.Session) {
	msg.handshake[s.Id()] = s.GetRemoteAddr()
	slog.Info("session:%d Open", s.Id())
}

func (msg *MsgServer) OnStop(s *frame.Session) {

}

// 找到玩家绑定的服务器
func getPlayerBindServer(player_id string) int {
	// panic("unimplemented")
	return consts.SERVER_TYPE_GAME
}
