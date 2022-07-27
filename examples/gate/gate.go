package gate

import "github.com/wspGreen/skyee/log"

var Gated = NewGate()

type Gate struct {
}

func NewGate() *Gate {
	return &Gate{}
}

func (g *Gate) init() {

}

func (g *Gate) OnClientRequest(protoname string) {
	// server_type := getProtoRouterServerType(protoname)
	// if server_type != g.server_type {
	// 	// 转发给其他服务器
	// 	g.forwardServerMessage()
	// } else {

	// }
}

// 接收服务间信息
func (g *Gate) OnServerMessage() {
	log.Println(" send OnServerMessage ok! ")
}

// 转发服务信息
func (g *Gate) forwardServerMessage() {
	log.Println(" send forwardServerMessage ok! ")
}

func (g *Gate) Move(a int) {
	log.Println(" send Move ", a)
}

func (g *Gate) Move1(a uint32) {
	log.Println(" send Move1 ", a)
}

func (g *Gate) Attck(name string) {
	log.Println(" send Attck ", name)
}
