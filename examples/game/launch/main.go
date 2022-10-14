package main

import (
	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/examples/game"
	"github.com/wspGreen/skyee/frame"
)

func main() {
	skyee.Start(func() {

		// skyee.SetFileLog("./log/gate.log")
		skyee.SetNats(game.GameBrokerHandler)
		skyee.NewService(game.Gamed, func(c *frame.SkyeeContext, params *skyee.OptionParam) {

			// skyee.SetWebSocket(c.GetId(), ":9101")

		})
	})
}
