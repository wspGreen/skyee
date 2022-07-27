package frame

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/wspGreen/skyee/component"
	"github.com/wspGreen/skyee/log"

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
	data []byte
}

func NewNetServerComp(addrtype string, addr string, ctxid uint32) *NetServerComp {
	return &NetServerComp{
		addr:      addr,
		addrtype:  addrtype,
		listener:  &http.Server{Addr: addr},
		forwardid: ctxid,
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
			net.listenWS(addr)
			// LogDebug("process listen end for msgque:%d", msgque.id)
		})
	}
	if addrtype == "http" {

		Go(func() {
			net.listenHttp(addr)
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

func (net *NetServerComp) listenHttp(url string) {
	mx := http.NewServeMux()
	mx.HandleFunc("/ClientPack", cors(func(w http.ResponseWriter, r *http.Request) {

		ctx := GetContext(net.forwardid)
		if ctx == nil {
			return
		}

		log.Println(fmt.Sprintf("[Webd:%d] Request Url:%s", ctx.id, r.RequestURI))

		msg := &Msg{}
		msg.Typename = "http"
		msg.CMD = "OnClientRequest"
		msg.Params = []interface{}{w, r}
		ctx.ac.Execute(msg)

	}))
	// http.ListenAndServe(":3002", mx)
	net.listener.Handler = mx
	Go(func() {
		log.Println("start http server port ", net.addr)
		net.listener.ListenAndServe()
	})

	// http.HandleFunc("/ClientPack", cors(func(w http.ResponseWriter, r *http.Request) {

	// 	ctx := GetContext(net.forwardid)
	// 	if ctx == nil {
	// 		return
	// 	}

	// 	log.Println(fmt.Sprintf("[Webd:%d] Request Url:%s", ctx.id, r.RequestURI))

	// 	msg := &Msg{}
	// 	msg.Typename = "http"
	// 	msg.CMD = "OnClientRequest"
	// 	msg.Params = []interface{}{w, r}
	// 	ctx.ac.Execute(msg)

	// }))

	// log.Println("start http server port ", net.addr)

	// net.listener.ListenAndServe()

}

func (net *NetServerComp) listenWS(url string) {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc(url, func(hw http.ResponseWriter, hr *http.Request) {
		c, err := upgrader.Upgrade(hw, hr, nil)
		if err != nil {
			// if stop == 0 && r.stop == 0 {
			log.Fatalf("accept failed  err:%v", err)
			// }
		} else {
			net.addSession(c)
			Go(func() {
				net.read(c)
			})

			Go(func() {
				net.write()
			})
		}
	})

	net.listener.ListenAndServe()
}

func (net *NetServerComp) addSession(con *websocket.Conn) {
	if con == nil {
		return
	}
	s := NewSession(con, net.GenSessionId())
	net.lock.Lock()
	net.sessions[s.id] = s
	net.lock.Unlock()
}

func (net *NetServerComp) GenSessionId() uint32 {
	return atomic.AddUint32(&net.seedSessionId, 1)
}

func (net *NetServerComp) read(conn *websocket.Conn) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}
		p := &Packet{data: data}
		net.onNetData(p)

	}
}

func (net *NetServerComp) onNetData(p *Packet) {
	// getContext(net.forwardid)

	// gonode.Send("socket", "OnClientRequest")
}

func (net *NetServerComp) write() {
	for {
		defer func() {
			if err := recover(); err != nil {
				// base.TraceCode(err)
			}
		}()

		// select {
		// case buff := <-s.sendChan:
		// 	if buff == nil { //信道关闭
		// 		return false
		// 	} else {
		// 		s.DoSend(buff)
		// 	}
		// }
	}
}
