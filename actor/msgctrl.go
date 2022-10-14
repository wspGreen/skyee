package actor

import (
	"sync/atomic"
)

type MsgCtrl struct {
	state uint32 // 0:无效 1:有效
	mark  chan bool
}

func NewMsgCtrl() *MsgCtrl {
	return &MsgCtrl{
		mark: make(chan bool, 1),
	}
}

func (c *MsgCtrl) Disable() {
	atomic.StoreUint32(&c.state, 0)
	// fmt.Println("Disable state  ", c.state, c.mark)

}

// 尝试设置状态为有效
func (c *MsgCtrl) trySetEnable() bool {
	b := atomic.LoadUint32(&c.state) == 0 && atomic.CompareAndSwapUint32(&c.state, 0, 1)
	// fmt.Println("trySetEnable state ", c.state, b, c.mark)
	return b
}

// 让消息可以继续执行
func (c *MsgCtrl) Enable() {
	if c.trySetEnable() {
		c.mark <- true
	}
}

func (c *MsgCtrl) IsEnable() chan bool {
	return c.mark
}
