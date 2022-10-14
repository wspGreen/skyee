package frame

import (
	"math"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/wspGreen/skyee/component"
	"github.com/wspGreen/skyee/slog"
	"github.com/wspGreen/skyee/snet"

	"github.com/gorilla/websocket"
)

type NetServerComp struct {
	component.Component
	addr          string
	addrtype      string
	listener      *http.Server
	seedSessionId uint32
	lock          *sync.RWMutex
	sessions      map[uint32]*Session
	forwardid     uint32
}

type Packet struct {
	Data []byte
}

func NewNetServerComp(addrtype string, addr string, ctxid uint32) *NetServerComp {
	return &NetServerComp{
		addr:      addr,
		addrtype:  addrtype,
		listener:  &http.Server{Addr: addr},
		forwardid: ctxid,
		lock:      &sync.RWMutex{},
		sessions:  make(map[uint32]*Session),
	}
}

func (net *NetServerComp) Start() {
	if net.IsInit {
		return
	}
	net.IsInit = true
	net.startServer(net.addrtype, net.addr)
}

func (net *NetServerComp) startServer(addrtype string, addr string) error {
	// addrs := strings.Split(addr, "://")
	if addrtype == "tcp" || addrtype == "all" {
		// listen, err := net.Listen("tcp", addrs[1])
		// if err == nil {
		// 	msgque := newTcpListen(listen, typ, handler, parser, addr)
		// 	Go(func() {
		// 		LogDebug("process listen for msgque:%d", msgque.id)
		// 		msgque.listen()
		// 		LogDebug("process listen end for msgque:%d", msgque.id)
		// 	})
		// } else {
		// 	LogError("listen on %s failed, errstr:%s", addr, err)
		// 	return err
		// }
	}
	if addrtype == "udp" || addrtype == "all" {
		// naddr, err := net.ResolveUDPAddr("udp", addrs[1])
		// if err != nil {
		// 	LogError("listen on %s failed, errstr:%s", addr, err)
		// 	return err
		// }
		// conn, err := net.ListenUDP("udp", naddr)
		// if err == nil {
		// 	msgque := newUdpListen(conn, typ, handler, parser, addr)
		// 	Go(func() {
		// 		LogDebug("process listen for msgque:%d", msgque.id)
		// 		msgque.listen()
		// 		LogDebug("process listen end for msgque:%d", msgque.id)
		// 	})
		// } else {
		// 	LogError("listen on %s failed, errstr:%s", addr, err)
		// 	return err
		// }
	}
	if addrtype == "ws" {
		// naddr := strings.SplitN(addrs[1], "/", 2)
		// url := "/"
		// if len(naddr) > 1 {
		// 	url = "/" + naddr[1]
		// }

		// msgque := newWsListen(naddr[0], url, MsgTypeCmd, handler, parser)
		Go(func() {
			// LogDebug("process listen for msgque:%d", msgque.id)
			net.listenWS()
			// LogDebug("process listen end for msgque:%d", msgque.id)
		})
	}
	if addrtype == "http" {

		Go(func() {
			net.listenHttp()
		})
	}
	return nil
}

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                                            // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", "true")                                                    //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                             //允许请求方法
		w.Header().Set("content-type", "application/json;charset=UTF-8")                                              //返回数据格式是json
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		f(w, r)
	}
}

func (net *NetServerComp) listenHttp() {
	mx := http.NewServeMux()
	mx.HandleFunc("/ClientPack", cors(func(w http.ResponseWriter, r *http.Request) {

		ctx := GetContext(net.forwardid)
		if ctx == nil {
			return
		}

		// log.Println(fmt.Sprintf("[Webd:%d] Request Url:%s", ctx.id, r.RequestURI))
		proto := ctx.GetProtoByType(PTYPE_SOCKET)
		if proto == nil {
			slog.Error("protoname is not exist : %v", PTYPE_SOCKET)
			return
		}

		// issue: http是否有必要用 proto.Pack
		msg := NewSkyeeMsg(PTYPE_SOCKET, 0, 0, proto.Pack([]interface{}{"OnClientRequest", w, r}))
		ctx.ac.Execute(msg)

	}))
	// http.ListenAndServe(":3002", mx)
	net.listener.Handler = mx
	Go(func() {
		slog.Info("start http server port %s", net.addr)
		err := net.listener.ListenAndServe()
		if err != nil {
			slog.Error("%v", err)
		}
	})

}

func (net *NetServerComp) listenWS() {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/", func(hw http.ResponseWriter, hr *http.Request) {
		// defer func() {
		// 	if err := recover(); err != nil {
		// 		slog.Stack()
		// 		slog.Error("error catch:%v", err)
		// 	}
		// }()
		c, err := upgrader.Upgrade(hw, hr, nil)
		if err != nil {
			if stop == 0 {
				slog.Fatal("accept failed  err:%v", err)
			} else {
				slog.Error("accept failed  err:%v", err)
			}
		} else {

			s := net.addSession(c)
			net.OnOpen(s)

			Go(func() {
				s.read()
			})

			Go(func() {
				s.write()
			})
		}
	})

	Go(func() {
		slog.Info("start ws server port %s", net.addr)
		err := net.listener.ListenAndServe()
		if err != nil {
			slog.Fatal("%v", err)
		}
	})
}

func (net *NetServerComp) GetSession(id uint32) *Session {
	net.lock.RLock()
	s := net.sessions[id]
	net.lock.RUnlock()
	return s
}

func (net *NetServerComp) addSession(con *websocket.Conn) *Session {
	if con == nil {
		return nil
	}
	s := NewSession(con, net.GenSessionId(), net)
	net.lock.Lock()
	net.sessions[s.id] = s
	net.lock.Unlock()

	return s
}

func (net *NetServerComp) delSession(id uint32) {
	net.lock.Lock()
	delete(net.sessions, id)
	net.lock.Unlock()
}

func (net *NetServerComp) GenSessionId() uint32 {
	if atomic.LoadUint32(&net.seedSessionId) == math.MaxUint32 {
		atomic.StoreUint32(&net.seedSessionId, 0)
	}

	return atomic.AddUint32(&net.seedSessionId, 1)
}

// func (net *NetServerComp) read(s *Session) {
// 	// buf := make([]byte, 2048)
// 	// for {
// 	// 	n, err := conn.Read(buf)
// 	// 	if err != nil {
// 	// 		// log.Println(fmt.Sprintf("Read message error: %s, session will be closed immediately", err.Error()))
// 	// 		return
// 	// 	}
// 	// }

// 	defer func() {
// 		if err := recover(); err != nil {
// 			// slog.Error("msgque read panic id:%v err:%v", r.id, err.(error))
// 			slog.Stack()
// 		}
// 		// r.Stop()
// 		// net.onStop(s)
// 	}()

// 	for {

// 		_, data, err := s.conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}

// 		// slog.Info("rev : %s", string(data))

// 		net.onNetData(data)

// 	}
// }

func (net *NetServerComp) OnOpen(s snet.ISession) {
	slog.Info("new session : %d, addr :%s", s.Id(), s.GetRemoteAddr())

	ctx := GetContext(net.forwardid)

	cmd := "OnOpen"
	if !ctx.ac.hasMethod(cmd) {
		return
	}

	proto := ctx.GetProtoByType(PTYPE_SOCKET)
	if proto == nil {
		slog.Error("protoname is not exist : %v", PTYPE_SOCKET)
		return
	}
	s = s.(*Session)
	params := proto.Pack([]interface{}{cmd, s})

	msg := NewSkyeeMsg(PTYPE_SOCKET, 0, 0, params)
	ctx.GetActor().Send(msg)
}

func (net *NetServerComp) OnNetData(s snet.ISession, data []byte) {

	ctx := GetContext(net.forwardid)
	if ctx == nil {
		slog.Error("not find context id %d", net.forwardid)
		return
	}
	proto := ctx.GetProtoByType(PTYPE_SOCKET)
	if proto == nil {
		slog.Error("protoname is not exist : %v", PTYPE_SOCKET)
		return
	}

	s = s.(*Session)
	params := proto.Pack([]interface{}{"OnClientRequest", s, data})
	msg := NewSkyeeMsg(PTYPE_SOCKET, 0, 0, params)
	ctx.GetActor().Send(msg)
	// gonode.Send("socket", "OnClientRequest")
}

// 每个seesion关闭时调用
func (net *NetServerComp) OnStop(s snet.ISession) {

	net.delSession(s.Id())
	slog.Info("stop session : %d", s.Id())

	// 通知绑定的actor关闭
	ctx := GetContext(net.forwardid)

	cmd := "OnStop"
	if !ctx.ac.hasMethod(cmd) {
		return
	}

	proto := ctx.GetProtoByType(PTYPE_SOCKET)
	if proto == nil {
		slog.Error("protoname is not exist : %v", PTYPE_SOCKET)
		return
	}
	s = s.(*Session)
	params := proto.Pack([]interface{}{cmd, s})

	msg := NewSkyeeMsg(PTYPE_SOCKET, 0, 0, params)
	ctx.GetActor().Send(msg)
}

func (net *NetServerComp) Stop() {

}
