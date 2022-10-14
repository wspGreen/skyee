package gate

import "github.com/wspGreen/skyee"

type BrokerHandler struct {
}

var GateBrokerHandler = NewGateBrokerHandler()

func NewGateBrokerHandler() *BrokerHandler {
	return &BrokerHandler{}
}

// 接收服务间信息
/*
	两种情况收到：
	1. 其他服务发给gate的包，在gate处理 (包来源 nats )
	2. 其他服务让gate转发client的包，只需转发给client （包来源client）
*/
func (g *BrokerHandler) OnServerMessage(actorid uint32, source uint32, cmd string, pid string, data []byte) {
	if actorid > 0 {
		skyee.Send(actorid, "cmd", cmd, pid, data)
	}
}
