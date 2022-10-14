package slib

import "github.com/wspGreen/skyee"

var DefaultBrokerHandler = NewDefaultBrokerHandler()

func NewDefaultBrokerHandler() *BrokerHandler {
	return &BrokerHandler{}
}

type BrokerHandler struct {
}

func (b *BrokerHandler) OnServerMessage(actorid uint32, source uint32, cmd string, pid string, data []byte) {
	if actorid > 0 {
		skyee.Send(actorid, "cmd", cmd, pid, data)
	}
}
