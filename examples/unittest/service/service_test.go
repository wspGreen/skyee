package service_test

import (
	"testing"
	"time"

	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/examples/unittest/service/servicehello"
	"github.com/wspGreen/skyee/slog"
)

func TestXxx(t *testing.T) {
	skyee.Start(func() {
		aid := skyee.NewService(servicehello.NewHello())
		slog.Info(aid)
		// for i := 0; i < 1; i++ {
		skyee.Send(aid, "cmd", "Move1", 222)
		// }

		skyee.Send(aid, "cmd", "Sleep")
		time.Sleep(time.Second)
		skyee.Send(aid, "cmd", "Move1", 222)
		// skyee.Send(aid, "cmd", "Move1", "rrrr")

		for {
			time.Sleep(time.Second * 1)
			skyee.Send(aid, "cmd", "Attck", "ssss")
		}

		// for i := 0; i < 2; i++ {
		// 	skyee.Send(aid, "CMD", "Move1", aid)
		// }

		// b := skyee.NewService(servicehello.Hello)

		// for i := 0; i < 10000; i++ {
		// 	skyee.NewService(servicehello.NewHello())
		// }

		// for i := 0; i < 2; i++ {
		// 	skyee.Send(b, "cmd", "Move1", uint32(i))
		// }

		// for i := 0; i < 10; i++ {
		// 	time.Sleep(1 * time.Second)
		// 	skyee.Send(b, "cmd", "Move1", uint32(i))
		// }

		// skyee.Send(b, "cmd", "Forward", aid)

	})

}
