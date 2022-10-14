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
		skyee.NewService(web.Web, func(c *frame.SkyeeContext, params *skyee.OptionParam) {
			skyee.SetHttp(c.GetId(), ":1002")
		})

		skyee.NewService(web.Web, func(c *frame.SkyeeContext, params *skyee.OptionParam) {
			skyee.SetHttp(c.GetId(), ":1004")
		})

	})

}
