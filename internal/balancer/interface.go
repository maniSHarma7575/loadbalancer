package balancer

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

type LB interface {
	Proxy(IncomingReq)
	Run()
	RunEventLoop()
	AddBackend(Backend)
	ChangeStrategy(string)
}

type Event interface {
	GetEventName() string
	GetData() interface{}
}
