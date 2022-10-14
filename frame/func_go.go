package frame

import (
	"sync/atomic"

	"github.com/wspGreen/skyee/slog"
)

func Go(fn func()) {
	waitAll.Add(1)
	var debugStr string

	id := atomic.AddUint64(&goid, 1)

	c := atomic.AddInt32(&gocount, 1)
	if slog.Log().Level() <= slog.LogLevelDebug {

		debugStr = slog.SimpleStack()
		slog.Debug("GOROUTINE START id:%d count:%d from:%s", id, c, debugStr)
	}
	go func() {
		// fn()
		Try(fn, nil)
		waitAll.Done()
		c = atomic.AddInt32(&gocount, -1)

		if slog.Log().Level() <= slog.LogLevelDebug {
			slog.Debug("GOROUTINE END id:%d count:%d from:%s", id, c, debugStr)
		}
	}()
}

func Go2(fn func(cstop chan struct{})) bool {
	if IsStopSvr() {
		return false
	}
	// waitAll.Add(1)
	var debugStr string

	id := atomic.AddUint64(&goid, 1)
	c := atomic.AddInt32(&gocount, 1)
	if slog.Log().Level() <= slog.LogLevelDebug {
		debugStr = slog.SimpleStack()
		slog.Debug("goroutine start id:%d count:%d from:%s", id, c, debugStr)
	}

	go func() {
		Try(func() { fn(stopChanForGo) }, nil)
		waitAll.Done()
		c = atomic.AddInt32(&gocount, -1)
		if slog.Log().Level() <= slog.LogLevelDebug {
			slog.Debug("goroutine end id:%d count:%d from:%s", id, c, debugStr)
		}
	}()
	return true
}

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			if handler == nil {
				slog.Stack()
				slog.Error("error catch:%v", err)
			} else {
				handler(err)
			}
			// atomic.AddInt32(&statis.PanicCount, 1)
			// statis.LastPanic = int(Timestamp)
		}
	}()
	fun()
}

func goForLog(fn func(cstop chan struct{})) bool {
	if IsStopSvr() {
		return false
	}
	waitAllForLog.Add(1)

	go func() {
		fn(stopChanForLog)
		waitAllForLog.Done()
	}()
	return true
}
