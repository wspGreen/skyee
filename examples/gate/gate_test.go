package gate_test

import (
	"testing"
	"time"

	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/examples/gate"
)

func TestXxx(t *testing.T) {
	id := skyee.NewService(gate.Gated)

	for i := 0; i < 1; i++ {
		skyee.Send(id, "CMD", "Move")
	}

	for i := 0; i < 2; i++ {
		skyee.Send(id, "CMD", "Move1", id)
	}

	b := skyee.NewService(gate.Gated)

	for i := 0; i < 2; i++ {
		skyee.Send(b, "CMD", "Move1", b)
	}

	for i := 0; i < 2; i++ {
		time.Sleep(1 * time.Second)
		skyee.Send(b, "CMD", "Move1", b)
	}

	skyee.WaitForSystemExit()
}
