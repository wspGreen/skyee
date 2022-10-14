package web_test

import (
	"testing"

	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/examples/web/web"
	"github.com/wspGreen/skyee/frame"
	"github.com/wspGreen/skyee/slog"
)

func TestXxx(t *testing.T) {

	skyee.Start(func() {
		slog.Log().SetLevel(slog.LogLevelInfo)

		skyee.NewService(web.Web, func(c *frame.SkyeeContext, params *skyee.OptionParam) {
			skyee.SetHttp(c.GetId(), ":1002")
		})

		skyee.NewService(web.Web, func(c *frame.SkyeeContext, params *skyee.OptionParam) {
			skyee.SetHttp(c.GetId(), ":1004")
		})

	})

}
