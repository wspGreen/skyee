package servicehello

import (
	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/log"
)

var Hello = NewHello()

type Hellod struct {
}

func NewHello() *Hellod {
	return &Hellod{}
}

func (g *Hellod) OnServerMessage() {
}

// func (g *Gate) Init() {

// }

func (g *Hellod) Move(a int) {
	log.Println(" send Move ", a)
}

func (g *Hellod) Move1(a uint32) {
	log.Println(" send Move1 ", a)
}

func (g *Hellod) Attck(name string) {
	log.Println(" send Attck ", name)
}

func (g *Hellod) Forward(id uint32) {
	log.Println(" send Forward to ", id)

	skyee.Send(id, "CMD", "Move", 123)
	skyee.Send(id, "CMD", "Move", 123)
}
