package gate_test

import (
	"testing"

	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/examples/gate"
	"github.com/wspGreen/skyee/frame"
)

func TestXxx(t *testing.T) {

	skyee.Start(func() {

		// skyee.SetFileLog("./log/gate.log")

		id := skyee.NewService(gate.NewGate(), func(c *frame.SkyeeContext, params *skyee.OptionParam) {

			// skyee.SetWebSocket(c.GetId(), ":9001")
		})
		// for i := 0; i < 3; i++ {
		// 	skyee.NewService(gate.NewGate(), func(c *frame.SkyeeContext, params *skyee.OptionParam) {

		// 	})

		// }

		uid := skyee.UniqueService(gate.NewGate())

		skyee.Send(uid, "cmd", "NewUniqueSvr", uint32(100))

		uid = skyee.UniqueService(gate.NewGate())
		skyee.Send(uid, "cmd", "NewUniqueSvr", uint32(101))

		skyee.Send("Gate", "cmd", "NewUniqueSvr", uint32(123))

		skyee.Send(id, "cmd", "Start")

		// for i := 0; i < 1; i++ {
		// 	skyee.Send(id, "CMD", "Move")
		// }

		// for i := 0; i < 2; i++ {
		// 	skyee.Send(id, "CMD", "Move1", id)
		// }

		// b := skyee.NewService(gate.Gated)

		// for i := 0; i < 2; i++ {
		// 	skyee.Send(b, "CMD", "Move1", b)
		// }

		// for i := 0; i < 10; i++ {
		// 	time.Sleep(1 * time.Second)
		// 	skyee.Send(b, "CMD", "Move1", b)
		// }

	})

}
