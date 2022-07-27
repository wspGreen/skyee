package service_test

import (
	"testing"
	"time"

	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/examples/unittest/service/servicehello"
)

func TestXxx(t *testing.T) {
	aid := skyee.NewService(servicehello.Hello)

	// for i := 0; i < 1; i++ {
	// 	skyee.Send(aid, "CMD", "Move")
	// }

	// for i := 0; i < 2; i++ {
	// 	skyee.Send(aid, "CMD", "Move1", aid)
	// }

	b := skyee.NewService(servicehello.Hello)

	for i := 0; i < 2; i++ {
		skyee.Send(b, "CMD", "Move1", uint32(i))
	}

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		skyee.Send(b, "CMD", "Move1", uint32(i))
	}

	skyee.Send(b, "CMD", "Forward", aid)

	skyee.WaitForSystemExit()
}
