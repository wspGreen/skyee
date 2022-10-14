package game_test

import (
	"testing"

	"github.com/wspGreen/skyee"
	"github.com/wspGreen/skyee/examples/game"
	"github.com/wspGreen/skyee/frame"
)

// type NatsHandler struct {
// }

// func (g *NatsHandler) GetSvrId() int {
// 	return consts.SERVER_TYPE_GAME
// }

// /*
// 	1. 收到gate转过来的包，属于客户端的包
// 	2. 收到其他服务发过来的包
// */
// func (g *NatsHandler) OnServerMessage(pid uint32, cmd string, data []byte) {
// 	slog.Info(" send OnServerMessage ok! ")

// 	///////////////////////////////////////
// 	// game
// 	// cmd := ""
// 	if getClientCmd(cmd) {
// 		onClientRequest(pid, data)
// 		return
// 	}

// 	fun := getServerCmd(cmd)
// 	fun()
// }

// // 处理客户端发来的包
// func onClientRequest(pid uint32, data []byte) {

// 	// 给对应agent处理？，在返回数据给gate
// 	Handle_XXX()
// }

// func getServerCmd(cmd string) func() {
// 	panic("unimplemented")
// }

// // 是否为客户端的包
// func getClientCmd(cmd string) bool {
// 	return true
// }

// func Handle_XXX() {
// 	send_response_by_gateway
// }

func TestXxx(t *testing.T) {
	skyee.Start(func() {

		// skyee.SetFileLog("./log/gate.log")
		// skyee.SetNats(&NatsHandler{})
		skyee.NewService(game.Gamed, func(c *frame.SkyeeContext, params *skyee.OptionParam) {

			// skyee.SetNats(game.Gamed)
			// skyee.SetWebSocket(c.GetId(), ":9101")

		})
	})
}
