package frame

import (
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/wspGreen/skyee/slog"
	"github.com/wspGreen/skyee/snet"
)

type Session struct {
	conn      *websocket.Conn
	id        uint32
	chwrite   chan *Packet
	stopstate int32
	handler   snet.INetHandler
}

func NewSession(con *websocket.Conn, id uint32, h snet.INetHandler) *Session {
	return &Session{
		id:      id,
		conn:    con,
		chwrite: make(chan *Packet, 64),
		handler: h,
	}
}

func (s *Session) Id() uint32 {
	return s.id
}

func (s *Session) GetRemoteAddr() string {
	return s.conn.RemoteAddr().String()
}

func (s *Session) read() {

	defer func() {
		if err := recover(); err != nil {
			// slog.Error("msgque read panic id:%v err:%v", r.id, err.(error))
			slog.Stack()
		}
		// r.Stop()
		// net.onStop(s)
		s.stop()
	}()

	for {
		s.updateStop()
		if s.isstop() {
			break
		}

		_, data, err := s.conn.ReadMessage()
		if err != nil {
			break
		}

		// slog.Debug("Client Request:%s", string(data))
		slog.Debug("Client Request: %.3f KB", float32(len(data))/1024)
		s.handler.OnNetData(s, data)
		// net.onNetData(data)

	}

}

func (s *Session) write() {
	defer func() {
		// if err := recover(); err != nil {
		// base.TraceCode(err)
		// }
		s.close()

		s.stop()
	}()
	var m *Packet
	// for !s.isstop() || m != nil {
	for {
		s.updateStop()
		if s.isstop() && m == nil {
			break
		}
		if m == nil {
			select {
			case m = <-s.chwrite:
			}

		}

		if m == nil {
			continue
		}

		err := s.Write(m.Data)
		if err != nil {
			break
		}

		m = nil
	}
}

func (s *Session) Write(data []byte) error {
	slog.Debug("Respone : %.3f KB", float32(len(data))/1024)
	err := s.conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		slog.Error("session write id:%v err:%v", s.id, err)
	}
	return err
}

func (s *Session) isstop() bool {
	return s.stopstate == 1
}

func (s *Session) updateStop() {
	if s.stopstate == 0 {
		if IsStopSvr() {
			s.stop()
		}
	}
}

// stopstate 设置为关闭状态，read write 根据这个状态跳出循环，关闭conn
func (s *Session) stop() {
	if atomic.CompareAndSwapInt32(&s.stopstate, 0, 1) {
		if s.chwrite != nil {
			close(s.chwrite)
		}
		s.handler.OnStop(s)
		// Go(func() {
		// 	if r.init {
		// 		r.handler.OnDelMsgQue(r)
		// 	}
		// 	r.available = false
		// 	r.baseStop()
		// })
	}
}

func (s *Session) Stop() {
	s.stop()
}

func (s *Session) close() {
	if s.conn != nil {
		s.conn.Close()
	}
}
