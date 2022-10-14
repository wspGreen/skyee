package game

import (
	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/slog"
)

type BrokerHandler struct {
}

var GameBrokerHandler = NewGameBrokerHandler()

func NewGameBrokerHandler() *BrokerHandler {
	return &BrokerHandler{}
}

/*
	1. 收到gate转过来的包，属于客户端的包
	2. 收到其他服务发过来的包
*/
func (g *BrokerHandler) OnServerMessage(actorid uint32, source uint32, cmd string, pid string, data []byte) {
	slog.Info("send OnServerMessage ok! ")
	if actorid > 0 {
		skyee.Send(actorid, "cmd", source, cmd, pid, data)
		return
	}
}
