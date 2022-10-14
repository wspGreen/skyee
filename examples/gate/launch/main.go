package main

import (
	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/examples/gate"
	"github.com/wspGreen/skyee/frame"
)

func main() {
	skyee.Start(func() {

		// skyee.SetFileLog("./log/gate.log")
		skyee.SetNats(gate.GateBrokerHandler)
		id := skyee.UniqueService(gate.NewGate(), func(c *frame.SkyeeContext, params *skyee.OptionParam) {
			skyee.SetWebSocket(c.GetId(), ":9001")

		})
		skyee.Send(id, "cmd", "Start")
		// skyee.Send(id, "cmd", "LoginHandler", "1001", "mark")

	})
}
