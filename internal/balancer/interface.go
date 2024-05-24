package balancer

import "net/http"

type Backend interface {
	Stringify() string
	GetHealthStatusUrl() string
	IncrementRequestCounter()
	UpdateIsHealthy(status bool)
	IsBackendHealthy() bool
	GetName() string
}

type IncomingReq interface {
	GetReqID() string
	GetHttpRequest() *http.Request
	GetHeadersAsMap() map[string]string
}

type BalancingStrategy interface {
	Init([]Backend)
	GetNextBackend(IncomingReq) Backend
	RegisterBackend(Backend)
	RefreshBackend(Backend)
	PrintTopology()
}

type LB interface {
	Run()
	RunEventLoop()
	AddBackend(Backend)
	ChangeStrategy(string)
	Observer
}

type Event interface {
	GetEventName() string
	GetData() interface{}
}

type Observer interface {
	BackendUp(backend Backend)
	BackendDown(backend Backend)
}

type Observable interface {
	Attach(observer Observer)
	Check()
}

type HealthCheckerInterface interface {
	Observable
	Start(interval int)
}
