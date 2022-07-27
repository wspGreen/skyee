package goread

import (
	"fmt"
	"testing"
	"time"
)

var ch chan struct{}

type taskInfo struct {
	id int
}

type Actor struct {
	active chan bool
	msgs   chan *taskInfo
}

func (a *Actor) Run() {
	// 目的:报错整个进程不终止
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("exe task error %v \n", err)
		}
	}()

	for msg := range a.msgs {
		fmt.Printf("通道收到任务 %d \n", msg.id)
		select {
		case <-a.active:
			a.Exe(msg)
		}
	}
	// for {
	// 	select {
	// 	case msg, ok := <-a.msgs:
	// 		if !ok {
	// 			fmt.Printf("通道收到任务错误 \n")
	// 		}
	// 		fmt.Printf("通道收到任务 %d \n", msg.id)
	// 		select {
	// 		case <-a.active:
	// 			a.Exe(msg)
	// 		}
	// 	}
	// }
}

func (a *Actor) Exe(t *taskInfo) {
	if t.id > 0 && t.id < 10 && 2%t.id == 0 {
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>执行任务 %d \n", t.id)
	a.active <- true
	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>完成任务 %d ，触发通道 \n", t.id)
}

func (a *Actor) AddTask(t *taskInfo) {
	a.msgs <- t
	// select {
	// case a.msgs <- t:
	// }

}

var startindex int

func inittask(a *Actor) {
	startindex = 10
	for i := 0; i < startindex; i++ {
		t := new(taskInfo)
		t.id = i
		fmt.Printf("初始化-任务 %d 加入 \n", t.id)
		a.AddTask(t)
	}
}

func TestXxx(t *testing.T) {
	ch = make(chan struct{})
	actor := Actor{
		active: make(chan bool, 1),
		msgs:   make(chan *taskInfo, 64),
	}
	go actor.Run()
	// 通道先放几个任务
	inittask(&actor)
	actor.active <- true

	// 开个go加任务，看任务的执行是否按顺序
	go func() {

		for i := startindex; i < startindex+10; i++ {
			t := new(taskInfo)
			t.id = i
			fmt.Printf("有任务 %d 加入 \n", t.id)
			actor.AddTask(t)
			time.Sleep(80 * time.Millisecond)
			// time.Sleep(1000 * time.Millisecond)

		}
	}()

	<-ch
}
