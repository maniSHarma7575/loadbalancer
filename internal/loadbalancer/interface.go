package loadbalancer

import "net"

type Backend interface {
	Stringify() string
	IncrementRequestCounter()
}

type IncomingReq interface {
	GetReqID() string
	GetSrcConn() net.Conn
}

type BalancingStrategy interface {
	Init([]Backend)
	GetNextBackend(IncomingReq) Backend
	RegisterBackend(Backend)
	PrintTopology()
}
