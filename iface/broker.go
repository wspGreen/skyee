package iface

type IBrokerHandler interface {
	OnServerMessage(actorid uint32, source uint32, cmd string, pid string, data []byte)
	// GetSvrId() int
}

type IBrokerPack interface {
	Pack(actorid uint32, params ...interface{}) ([]byte, error)

	// params = {actorid, cmd, pid, data}
	UnPack(ndata []byte) (params []interface{}, err error)
}
