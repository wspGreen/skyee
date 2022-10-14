package slib

import (
	"strconv"

	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/slib/consts"
	"github.com/wspGreen/skyee/slog"
)

// skyee 自带的服务 ，也可以自己实现

// type BaseServer struct {
// }

// // 转发服务信息
// func (b *BaseServer) ForwardServerMessage(srvid string, pid uint32, data []byte) {
// 	slog.Info(" send forwardServerMessage ok! ")
// 	npkt := skyee.Nats().Pack(pid, data)
// 	if srvid != "" {
// 		skyee.Nats().SendToServer(srvid, npkt)
// 	} else {
// 		skyee.Nats().SendToGroupServer()
// 	}
// }

// func (b *BaseServer) SendResponseByGate() {

// }

// 转发服务信息
func ForwardServerMessage(srvid string, dest uint32, source uint32, cmd string, pid string, data []byte) {
	// npkt, _ := skyee.Nats().Pack(broker.NewBrokerMessage(pid, 0, data))
	if srvid != "" {
		// skyee.Nats().SendToServer(srvid, npkt)
		skyee.Nats().Send(srvid, dest, source, cmd, pid, data)
	} else {
		skyee.Nats().SendToGroupServer()
	}
}

func SendResponseByGate(pid string, dest uint32, data []byte) {
	slog.Debug("Server To Gate  (playerid:%d,gateid:%d)", pid, consts.SERVER_TYPE_GATE)
	ForwardServerMessage(strconv.Itoa(consts.SERVER_TYPE_GATE), dest, 0, "OnServerMessage", pid, data)
}
