package slogimp

import (
	"fmt"
	"os"
)

// 计划把 frame.log 的代码搬到这里

type ConsoleLogger struct {
	Ln bool
}

func (r *ConsoleLogger) Write(str string) {
	if r.Ln {
		fmt.Println(str)
	} else {
		fmt.Print(str)
	}
}

type OnFileLogFull func(path string)
type OnFileLogTimeout func(path string) int
type FileLogger struct {
	Path        string
	Ln          bool
	DatePostfix bool
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
