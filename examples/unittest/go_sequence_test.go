package unittest_test

import (
	"fmt"
	"testing"
	"time"
)

var isactive chan bool
var ch chan struct{}

var count int

// 串行
func Work() {
	fmt.Printf("Go 串行--------------------------------------------------\n")
	for i := 0; i < count; i++ {
		select {
		case job := <-isactive:
			go task(i, job)
		}
	}
}

func task(n int, b bool) {
	defer func() {
		isactive <- true
	}()
	// size := randTime()
	if n == 5 {
		call(n)

	} else {
		fmt.Printf("worker %2d process \n", n)
	}

}

func call(n int) time.Duration {
	size := time.Duration(1000)
	time.Sleep(size * time.Millisecond)
	fmt.Printf("worker %2d process , time %dms\n", n, size)
	return size
}

// 并行
func Work02() {
	fmt.Printf("Go 并行--------------------------------------------------\n")
	for i := 0; i < count; i++ {
		go task02(i, true)

	}
}

func task02(n int, b bool) {

	size := time.Duration(100)
	time.Sleep(size * time.Millisecond)
	fmt.Printf("worker %2d process , time %dms\n", n, size)
}

///////////////////////////
func Work03() {
	fmt.Printf("Go --------------------------------------------------\n")
	for i := 0; i < count; i++ {
		select {
		case job := <-isactive:
			go task03(i, job)
		}
	}
}

func task03(n int, b bool) {

	// size := randTime()
	if n == 5 || n == 6 {
		isactive <- true
		call01(n)

	} else {
		fmt.Printf("非阻塞函数 %2d 执行完成 \n", n)
		isactive <- true
	}

}

func call01(n int) time.Duration {
	fmt.Printf("阻塞函数 %2d 执行开始 \n", n)
	size := time.Duration(1000)
	time.Sleep(size * time.Millisecond)
	fmt.Printf("阻塞函数 %2d 执行完成 , time %dms\n", n, size)
	return size
}

func TestStart(t *testing.T) {
	count = 10
	isactive = make(chan bool)
	ch = make(chan struct{}, 1)

	go Work03()

	time.Sleep(500 * time.Millisecond)
	fmt.Printf("开始激活\n")
	isactive <- true

	// Work02()

	<-ch
}
