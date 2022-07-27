package frame

import "github.com/gorilla/websocket"

type Session struct {
	conn *websocket.Conn
	id   uint32
}

func NewSession(con *websocket.Conn, id uint32) *Session {
	return &Session{id: id, conn: con}
}
