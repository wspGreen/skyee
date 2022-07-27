package main

import (
	"fmt"
	"time"

	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/examples/gate"
	"github.com/wspGreen/skyee/frame"
)

func main() {
	// args := os.Args
	// var a uint32
	a := skyee.NewService(gate.Gated, func(ctx *frame.SkyeeContext, opt *skyee.OptionParam) {

		// gonode.SetWebSocket(ctx.GetId())
	})

	for i := 0; i < 1; i++ {
		skyee.Send(a, "CMD", "Move")
	}

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		skyee.Send(a, "CMD", "Move", i)
	}

	fmt.Printf(">>>>>>>>>>>>> DONE")

	skyee.WaitForSystemExit()

}

func init() {
	fmt.Printf(">>>>>>>>>>>>> INIT")
}
