package lib

import "fmt"

type BalancingStrategy interface {
	Init([]*Backend)
	GetNextBackend(IncomingReq) *Backend
	RegisterBackend(*Backend)
	PrintTopology()
}

type RoundRobinBS struct {
	Index    int
	Backends []*Backend
}

type StaticBS struct {
	Index    int
	Backends []*Backend
}

func (rrbs *RoundRobinBS) Init(backends []*Backend) {
	rrbs.Index = 0
	rrbs.Backends = backends
}

func (rrbs *RoundRobinBS) GetNextBackend(IncomingReq) *Backend {
	rrbs.Index = (rrbs.Index + 1) % len(rrbs.Backends)
	return rrbs.Backends[rrbs.Index]
}

func (rrbs *RoundRobinBS) RegisterBackend(backend *Backend) {
	rrbs.Backends = append(rrbs.Backends, backend)
}

func (rrbs *RoundRobinBS) PrintTopology() {
	for index, backend := range rrbs.Backends {
		fmt.Println(fmt.Sprintf("		[%d] %s", index, backend))
	}
}

func NewRoundRobinBS(backends []*Backend) *RoundRobinBS {
	strategy := new(RoundRobinBS)
	strategy.Init(backends)
	return strategy
}

func (sbs *StaticBS) Init(backends []*Backend) {
	sbs.Index = 0
	sbs.Backends = backends
}

func (sbs *StaticBS) GetNextBackend(IncomingReq) *Backend {
	return sbs.Backends[sbs.Index]
}

func (sbs *StaticBS) RegisterBackend(backend *Backend) {
	sbs.Backends = append(sbs.Backends, backend)
}

func (sbs *StaticBS) PrintTopology() {
	for index, backend := range sbs.Backends {
		if index == sbs.Index {
			fmt.Println(fmt.Sprintf("	[%s] %s", "x", backend))
		} else {
			fmt.Println(fmt.Sprintf(" [%s] %s", " ", backend))
		}
	}
}

func NewStaticBS(backends []*Backend) *StaticBS {
	strategy := new(StaticBS)
	strategy.Init(backends)
	return strategy
}
