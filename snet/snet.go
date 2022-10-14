package snet

type ISession interface {
	Id() uint32
	GetRemoteAddr() string
}

type INetHandler interface {
	OnOpen(s ISession)
	OnNetData(s ISession, data []byte)

	OnStop(s ISession)
}
