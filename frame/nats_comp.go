package frame

import (
	"github.com/nats-io/nats.go"
	"github.com/wspGreen/skyee/broker"
	"github.com/wspGreen/skyee/component"
	"github.com/wspGreen/skyee/iface"
	"github.com/wspGreen/skyee/slog"
)

type NatsComp struct {
	component.Component
	handler iface.IBrokerHandler
	pack    iface.IBrokerPack
	nc      *nats.Conn
}

func NewNatsComp(h iface.IBrokerHandler) *NatsComp {
	n := &NatsComp{
		handler: h,
		pack:    broker.DefaultBrokerPack}

	return n
}

func (n *NatsComp) Start() {
	if n.IsInit {
		return
	}
	n.IsInit = true

	nc, err := nats.Connect(broker.NATS_ADDR,
		nats.DisconnectHandler(func(_ *nats.Conn) { slog.Error("disconnected from nats!") }),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			slog.Warn("reconnected to nats server %s with address %s in cluster %s!", nc.ConnectedServerId(), nc.ConnectedAddr(), nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			err := nc.LastError()
			if err == nil {
				slog.Info("nats connection closed with no error.")
				return
			}

			slog.Error("nats connection closed. reason: %q", nc.LastError())
			// if appDieChan != nil {
			// 	appDieChan <- true
			// }
		}),
	)
	if err != nil {
		slog.Error(err)
		return
	}
	slog.Info("nats connection succeed %s", broker.NATS_ADDR)
	n.nc = nc

}
func (n *NatsComp) Subscribe(subj string) {
	n.nc.Subscribe(subj, n.on_nats_message)
}

func (n *NatsComp) on_nats_message(msg *nats.Msg) {
	// slog.Info("on_nats_message : %s", string(msg.Data))

	params, err := n.pack.UnPack(msg.Data)
	if err != nil {
		slog.Error(err)
	}

	// 应该要放actor处理
	actorid := params[0].(uint32)
	source := params[1].(uint32)
	cmd := params[2].(string)
	pid := params[3].(string)
	data := params[4].([]byte)
	n.handler.OnServerMessage(actorid, source, cmd, pid, data)
}

func (n *NatsComp) SendToServer(svrid string, data []byte) {
	// if n.nc == nil {
	// 	return
	// }

	// if err := n.nc.Publish(svrid, data); err != nil {
	// 	slog.Error("%v", err)
	// 	return
	// }
}

func (n *NatsComp) SendToGroupServer() {

}

func (n *NatsComp) SendToAllServer() {

}

func (n *NatsComp) Pack(actorid uint32, params ...interface{}) ([]byte, error) {
	return n.pack.Pack(actorid, params...)
}

func (n *NatsComp) SetPack(pack iface.IBrokerPack) {
	n.pack = pack
}

func (n *NatsComp) Send(node string, dest uint32, params ...interface{}) {
	if n.nc == nil {
		return
	}
	// npkt, _ := n.pack.Pack(broker.NewBrokerMessage(pid, 0, data))
	npkt, _ := n.pack.Pack(dest, params...)
	if err := n.nc.Publish(node, npkt); err != nil {
		slog.Error("%v", err)
		return
	}
}
