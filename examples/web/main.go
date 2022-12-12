package main

import (
	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/examples/web/web"
	"github.com/wspGreen/skyee/frame"
	"github.com/wspGreen/skyee/slog"
)

func main() {

	skyee.Start(func() {
		skyee.SetFileLog("./weblog/web.log").SetLevel(slog.LogLevelInfo)
		a := web.NewWeb()
		// slog.Error("a addr : %p", a)
		skyee.NewService(a, func(c *frame.SkyeeContext, params *skyee.OptionParam) {
			skyee.SetHttp(c.GetId(), ":1002")
		})

		b := web.NewWeb()
		// slog.Error("b addr : %p", b)
		skyee.NewService(b, func(c *frame.SkyeeContext, params *skyee.OptionParam) {
			skyee.SetHttp(c.GetId(), ":1004")
		})

	})

}
