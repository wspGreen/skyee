package frame

import (
	"time"

	"github.com/wspGreen/skyee/slog"
)

func ParseTime(str string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", str)
}

func Date() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Date2() string {
	return time.Now().Format("2006-01-02 15_04_05")
}

func Date3() string {
	return time.Now().Format("2006-01-02")
}

func UnixTime(sec, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}

func UnixMs() int64 {
	return time.Now().UnixNano() / 1000000
}

func Now() time.Time {
	return time.Now()
}

func NewTimer(ms int) *time.Timer {
	return time.NewTimer(time.Millisecond * time.Duration(ms))
}

func NewTicker(ms int) *time.Ticker {
	return time.NewTicker(time.Millisecond * time.Duration(ms))
}

func After(ms int) <-chan time.Time {
	return time.After(time.Millisecond * time.Duration(ms))
}

func Tick(ms int) <-chan time.Time {
	return time.Tick(time.Millisecond * time.Duration(ms))
}

func Sleep(ms int) {
	time.Sleep(time.Millisecond * time.Duration(ms))
}

func SetTimeout(inteval int, fn func(...interface{}) int, args ...interface{}) {
	if inteval < 0 {
		slog.Error("new timerout inteval:%v", inteval)
		return
	}
	slog.Info("new timerout inteval:%v", inteval)

	Go2(func(cstop chan struct{}) {
		var tick *time.Timer
		for inteval > 0 {
			tick = time.NewTimer(time.Millisecond * time.Duration(inteval))
			select {
			case <-cstop:
				inteval = 0
			case <-tick.C:
				tick.Stop()
				inteval = fn(args...)
			}
		}
		if tick != nil {
			tick.Stop()
		}
	})
}
