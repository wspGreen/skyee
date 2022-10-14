package log_test

import (
	"testing"
	"time"

	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/frame"
	"github.com/wspGreen/skyee/slog"
)

func run_log1() {
	frame.Go(func() {
		for {
			slog.Info("this is log info")
			time.Sleep(1 * time.Second)

		}
	})
}

func run_log2() {
	// s := path.Dir("D:/Project/mykingdom/server/skyee/examples")
	// fmt.Println(s)
	// 10 秒后保存一次日志
	slog.SetLogger(&frame.FileLogger{Path: "./log/logtest.log", Ln: true, Timeout: 10,
		OnTimeout: func(path string) int {
			// 可用来上传文件
			slog.Info(path)
			return 10
		}})
	frame.Go(func() {
		for {
			slog.Info("this is filelog2 info %v", frame.Date2())
			time.Sleep(1 * time.Second)

		}
	})
}

func run_log3() {
	skyee.SetFileLog("./log/game.log")
	// frame.DefLog.SetLogger(&frame.FileLogger{Path: "./log/game.log", Ln: true})
	frame.Go(func() {
		for {
			slog.Info("this is filelog info")
			slog.Info("this is filelog info1 %v", frame.Date2())
			time.Sleep(1 * time.Second)

		}
	})
}

func run_log_full() {

	slog.SetLogger(&frame.FileLogger{Path: "./log/logtest.log", Ln: true, MaxSize: 3 * 1024,
		OnFull: func(path string) {

			slog.Info("file full : %v", path)
		}})
	frame.Go(func() {
		for {
			slog.Info("this is filelog full info %v", frame.Date2())
			time.Sleep(1 * time.Second)

		}
	})
}

func TestXxx(t *testing.T) {
	skyee.Start(func() {
		run_log_full()

	})
}
