package unittest_test

import (
	"sync"
	"testing"
)

// 此 dome 为 map 并发读写报错范例

var list = make(map[int]*Data)

var lock sync.Mutex
var lockrw sync.RWMutex

type Data struct {
	num int
}

func TestSta(t *testing.T) {
	list[1] = &Data{1}
	go func() {
		for {
			lockrw.Lock()
			list[2] = &Data{2}
			lockrw.Unlock()
		}
	}()

	go func() {
		for {
			// lockrw.RLock()
			_ = list[1]   // 非线程安全
			_ = len(list) // 线程安全
			// lockrw.RUnlock()

		}
	}()

	select {}

}
