package frame

import (
	"os"
	"sync"
	"sync/atomic"

	"github.com/wspGreen/skyee/iface"
	"github.com/wspGreen/skyee/slog"
	"github.com/wspGreen/skyee/slog/slogimp"
)

var PTYPE_SOCKET = uint8(1)
var PTYPE_CMD = uint8(2)
var PTYPE_RESPONSE = uint8(3)

var stopChanForSys = make(chan os.Signal, 1)
var stopChanForGo = make(chan struct{})
var stopChanForLog = make(chan struct{}) // 通知log关闭

var stop int32       //停止状态
var stopForLog int32 //log停止状态

var gocount int32 //goroutine数量
var goid uint64

type Service func() iface.IHander

type WaitGroup struct {
	count int64
}

func (r *WaitGroup) Add(delta int) {
	atomic.AddInt64(&r.count, int64(delta))
}

func (r *WaitGroup) Done() {
	atomic.AddInt64(&r.count, -1)
}

func (r *WaitGroup) Wait() {
	for atomic.LoadInt64(&r.count) > 0 {
		Sleep(1)
	}
}

func (r *WaitGroup) TryWait() bool {
	return atomic.LoadInt64(&r.count) == 0
}

var waitAll = &WaitGroup{} //等待所有goroutine
var waitAllForLog sync.WaitGroup

// var DefLog *Log //日志

// var DefLog slog.ISLog

func init() {
	// runtime.GOMAXPROCS(runtime.NumCPU())
	// DefLog = NewLog(10000)
	slog.SetLog(NewLog(10000))
	// DefLog = slog.DefLog
	slog.Log().SetLogger(&slogimp.ConsoleLogger{Ln: true})
}
