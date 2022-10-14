package frame

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/wspGreen/skyee/slog"
)

// type ILogger interface {
// 	Write(str string)
// }

// type ConsoleLogger struct {
// 	Ln bool
// }

// func (r *ConsoleLogger) Write(str string) {
// 	if r.Ln {
// 		fmt.Println(str)
// 	} else {
// 		fmt.Print(str)
// 	}
// }

type OnFileLogFull func(path string)
type OnFileLogTimeout func(path string) int
type FileLogger struct {
	Path string
	Ln   bool
	// DatePostfix bool
	/*
		Timeout后Path的文件被改名,并创建新的Path文件继续写入
	*/
	Timeout   int //0表示不设置, 单位s
	MaxSize   int //0表示不限制，最大大小
	OnFull    OnFileLogFull
	OnTimeout OnFileLogTimeout // timeout后回调

	size     int
	file     *os.File
	filename string
	extname  string
	dirname  string
}

func (r *FileLogger) Write(str string) {
	if r.file == nil {
		return
	}

	newsize := r.size
	if r.Ln {
		newsize += len(str) + 1
	} else {
		newsize += len(str)
	}

	if r.MaxSize > 0 && newsize >= r.MaxSize {
		r.file.Close()
		r.file = nil
		newpath := r.dirname + "/" + r.filename + fmt.Sprintf("_%v", Date2()) + r.extname
		os.Rename(r.Path, newpath)
		file, err := os.OpenFile(r.Path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if err == nil {
			r.file = file
		}
		r.size = 0
		if r.OnFull != nil {
			r.OnFull(newpath)
		}
	}

	if r.file == nil {
		return
	}

	if r.Ln {
		r.file.WriteString(str)
		r.file.WriteString("\n")
		r.size += len(str) + 1
	} else {
		r.file.WriteString(str)
		r.size += len(str)
	}
}

// type LogLevel int

// const (
// 	LogLevelAllOn  slog.LogLevel = iota //开放说有日志
// 	LogLevelDebug                       //调试信息
// 	LogLevelInfo                        //资讯讯息
// 	LogLevelWarn                        //警告状况发生
// 	LogLevelError                       //一般错误，可能导致功能不正常
// 	LogLevelFatal                       //严重错误，会导致进程退出
// 	LogLevelAllOff                      //关闭所有日志
// )

type Log struct {
	// slog.ISLog
	logger         [32]slog.ILogger
	cwrite         chan string
	ctimeout       chan *FileLogger
	bufsize        int
	stop           int32
	preLoggerCount int32
	loggerCount    int32
	level          slog.LogLevel
}

func (r *Log) initFileLogger(f *FileLogger) *FileLogger {
	if f.file == nil {
		f.Path, _ = filepath.Abs(f.Path)
		// filepath.Dir()
		f.dirname = filepath.Dir(f.Path)
		f.extname = path.Ext(f.Path)
		f.filename = filepath.Base(f.Path[0 : len(f.Path)-len(f.extname)])
		// // game.log = game_0-0-0.log
		// if f.DatePostfix {
		// 	f.Path = f.dirname + "/" + f.filename + fmt.Sprintf("_%v", Date3()) + f.extname
		// }
		os.MkdirAll(f.dirname, 0777)
		file, err := os.OpenFile(f.Path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if err == nil {
			f.file = file
			info, err := f.file.Stat()
			if err != nil {
				return nil
			}

			f.size = int(info.Size())
			if f.Timeout > 0 {
				SetTimeout(f.Timeout*1000, func(...interface{}) int {
					defer func() { recover() }()
					r.ctimeout <- f
					return 0
				})
			}

			return f
		} else {
			log.Fatalf("open file err: %v", err)
		}
	}
	return nil
}

func (r *Log) start() {
	goForLog(func(cstop chan struct{}) {
		var i int32
		for !r.IsStop() {
			select {
			case s, ok := <-r.cwrite:
				if ok {
					for i = 0; i < r.loggerCount; i++ {
						r.logger[i].Write(s)
					}
				}
			case c, ok := <-r.ctimeout:
				if ok {
					c.file.Close()
					c.file = nil
					newpath := c.dirname + "/" + c.filename + fmt.Sprintf("_%v", Date2()) + c.extname
					e := os.Rename(c.Path, newpath)
					if e != nil {
						slog.Info("rename error: %v", e)
					}
					file, err := os.OpenFile(c.Path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
					if err == nil {
						c.file = file
					}
					c.size = 0
					if c.OnTimeout != nil {
						nc := c.OnTimeout(newpath)
						if nc > 0 {
							SetTimeout(nc*1000, func(...interface{}) int {
								defer func() { recover() }()
								r.ctimeout <- c
								return 0
							})
						}
					}
				}
			case <-cstop:
			}
		}

		for s := range r.cwrite {
			for i = 0; i < r.loggerCount; i++ {
				r.logger[i].Write(s)
			}
		}
	})
}

func (r *Log) Stop() {
	if atomic.CompareAndSwapInt32(&r.stop, 0, 1) {
		close(r.cwrite)
		close(r.ctimeout)
	}
}

func (r *Log) SetLogger(logger slog.ILogger) bool {
	if r.preLoggerCount >= 31 {
		return false
	}
	if f, ok := logger.(*FileLogger); ok {
		if r.initFileLogger(f) == nil {
			return false
		}
	}
	r.logger[atomic.AddInt32(&r.preLoggerCount, 1)] = logger
	atomic.AddInt32(&r.loggerCount, 1)
	return true

}
func (r *Log) Level() slog.LogLevel {
	return r.level
}
func (r *Log) SetLevel(level slog.LogLevel) {
	r.level = level
}

func isLogStop() bool {
	return stopForLog == 1
}

func (r *Log) IsStop() bool {
	if r.stop == 0 {
		if isLogStop() {
			r.Stop()
		}
	}
	return r.stop == 1
}

func (r *Log) write(levstr string, v ...interface{}) {
	defer func() { recover() }()
	if r.IsStop() {
		return
	}
	prefix := levstr
	_, file, line, ok := runtime.Caller(3)
	if ok {
		i := strings.LastIndex(file, "/") + 1
		prefix = fmt.Sprintf("[%s][%s][%s:%d]:", levstr, Date(), (string)(([]byte(file))[i:]), line)
	}
	if len(v) > 1 {
		r.cwrite <- prefix + fmt.Sprintf(v[0].(string), v[1:]...)
	} else {
		r.cwrite <- prefix + fmt.Sprint(v[0])
	}
}

func (r *Log) Debug(v ...interface{}) {
	if r.level <= slog.LogLevelDebug {
		r.write("D", v...)
	}
}

func (r *Log) Info(v ...interface{}) {
	if r.level <= slog.LogLevelInfo {
		r.write("I", v...)
	}
}

func (r *Log) Warn(v ...interface{}) {
	if r.level <= slog.LogLevelWarn {
		r.write("W", v...)
	}
}

func (r *Log) Error(v ...interface{}) {
	if r.level <= slog.LogLevelError {
		r.write("E", v...)
	}
}

func (r *Log) Fatal(v ...interface{}) {
	if r.level <= slog.LogLevelFatal {
		r.write("FATAL", v...)
	}
}

func (r *Log) Write(v ...interface{}) {
	defer func() { recover() }()
	if r.IsStop() {
		return
	}

	if len(v) > 1 {
		r.cwrite <- fmt.Sprintf(v[0].(string), v[1:]...)
	} else if len(v) > 0 {
		r.cwrite <- fmt.Sprint(v[0])
	}
}

func NewLog(bufsize int, logger ...slog.ILogger) *Log {
	log := &Log{
		bufsize:        bufsize,
		cwrite:         make(chan string, bufsize),
		ctimeout:       make(chan *FileLogger, 32),
		level:          slog.LogLevelAllOn,
		preLoggerCount: -1,
	}
	for _, l := range logger {
		log.SetLogger(l)
	}
	log.start()
	return log
}

// func LogInfo(v ...interface{}) {
// 	DefLog.Info(v...)
// }

// func LogDebug(v ...interface{}) {
// 	DefLog.Debug(v...)
// }

// func LogError(v ...interface{}) {
// 	DefLog.Error(v...)
// }

// func LogFatal(v ...interface{}) {
// 	DefLog.Fatal(v...)
// }

// func LogWarn(v ...interface{}) {
// 	DefLog.Warn(v...)
// }

// func LogStack() {
// 	buf := make([]byte, 1<<12)
// 	slog.Error(string(buf[:runtime.Stack(buf, false)]))
// }

// func LogSimpleStack() string {
// 	_, file, line, _ := runtime.Caller(2)
// 	i := strings.LastIndex(file, "/") + 1
// 	i = strings.LastIndex((string)(([]byte(file))[:i-1]), "/") + 1

// 	return Sprintf("%s:%d", (string)(([]byte(file))[i:]), line)
// }
