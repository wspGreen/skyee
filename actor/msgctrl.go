package actor

import (
	"sync/atomic"
)

type MsgCtrl struct {
	state uint32
	mark  chan bool
}

func NewMsgCtrl() *MsgCtrl {
	return &MsgCtrl{
		mark: make(chan bool, 1),
	}
}

func (c *MsgCtrl) EnableMarkState() {
	// fmt.Println("Enable  ", c.mark)
	atomic.StoreUint32(&c.state, 0)
}

func (c *MsgCtrl) CheckAndDisable() bool {
	b := atomic.LoadUint32(&c.state) == 0 && atomic.CompareAndSwapUint32(&c.state, 0, 1)
	// fmt.Println("CheckAndDisable  ", b, c.mark)
	return b
}

func (c *MsgCtrl) Mark() {
	if c.CheckAndDisable() {
		c.mark <- true
	}
}

func (c *MsgCtrl) GetMark() chan bool {
	return c.mark
}
