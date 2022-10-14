package broker_test

import (
	"fmt"
	"testing"

	"github.com/wspGreen/skyee/broker"
)

func TestXxx(t *testing.T) {
	p := broker.DefaultBrokerPack
	// msg := broker.NewBrokerMessage(666, []byte("88888"))
	b, _ := p.Pack(0, "")

	msg1, _ := p.UnPack(b)
	fmt.Printf(" msg : %v \n", msg1)
}
