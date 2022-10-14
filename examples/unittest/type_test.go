package unittest_test

import (
	"fmt"
	"testing"
)

type SER func()

func GameD() {
	fmt.Println("GameD")
}

func NewSer(s SER) {
	s()
}

var ser SER = GameD

func Test_T(t *testing.T) {
	NewSer(ser)
}
