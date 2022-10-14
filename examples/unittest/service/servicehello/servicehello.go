package servicehello

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/iface"
	"github.com/wspGreen/skyee/slog"
)

var Hello = NewHello()

type Hellod struct {
}

func NewHello() *Hellod {
	return &Hellod{}
}

func (g *Hellod) Init(a iface.IActor) {
	skyee.Dispatch(
		a.GetId(),
		"cmd",
		func(session uint32, source uint32, cmd string, params []interface{}) []interface{} {
			ret := a.FunCall(cmd, params)
			// paramList = append(paramList, reflect.ValueOf(session))
			// ret := m.Func.Call(paramList)
			if ret != nil {

			}
			return nil
		},
	)
	// id := 1
	// dispatch := func(m *reflect.Method, paramList []reflect.Value) bool {
	// 	ret := m.Func.Call(paramList)
	// 	if ret != nil {

	// 	}
	// 	return true
	// }
}

func (g *Hellod) OnServerMessage() {

}

// func (g *Gate) Init() {

// }

func (g *Hellod) Move(session uint32, a int) {
	slog.Info(" send Move %v session:%v", a, session)
}

func (g *Hellod) Move1(a int) {
	slog.Info(" send Move1 %v", a)

	fmt.Println(Colorize("rrrrrr", FgHiRed))
}

func (g *Hellod) Attck(name string) {
	slog.Info(" send Attck %v", name)
}

func (g *Hellod) Forward(id uint32) {
	slog.Info(" send Forward to %v", id)

	skyee.Send(id, "cmd", "Move", 123)
	skyee.Send(id, "cmd", "Move", 123)
}

func (g *Hellod) Sleep() {
	fmt.Println("Sleep start ", time.Now())
	time.Sleep(time.Second * 10)
	fmt.Println("Sleep end ", time.Now())
}

// Color defines a single SGR Code
type Color int

// Foreground text colors
const (
	FgBlack Color = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack Color = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Colorize a string based on given color.
func Colorize(s string, c Color) string {
	return fmt.Sprintf("\033[1;%s;40m%s\033[0m", strconv.Itoa(int(c)), s)
}
