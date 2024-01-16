package strategy

import (
	"fmt"

	loadbalancer "github.com/maniSHarma7575/loadbalancer/internal/balancer"
)

type StaticBS struct {
	Index    int
	Backends []loadbalancer.Backend
}

func (sbs *StaticBS) Init(backends []loadbalancer.Backend) {
	sbs.Index = 0
	sbs.Backends = backends
}

func (sbs *StaticBS) GetNextBackend(loadbalancer.IncomingReq) loadbalancer.Backend {
	return sbs.Backends[sbs.Index]
}

func (sbs *StaticBS) RefreshBackend(backend loadbalancer.Backend) {
	idx := FindBackendIndex(sbs.Backends, backend)

	if backend.IsBackendHealthy() && idx == -1 {
		sbs.Backends = append(sbs.Backends, backend)
	} else if idx != -1 && !backend.IsBackendHealthy() {
		sbs.Backends = append(sbs.Backends[:idx], sbs.Backends[idx+1:]...)
	}
}

func (sbs *StaticBS) RegisterBackend(backend loadbalancer.Backend) {
	sbs.Backends = append(sbs.Backends, backend)
}

func (sbs *StaticBS) PrintTopology() {
	for index, backend := range sbs.Backends {
		if index == sbs.Index {
			fmt.Printf("	[%s] %s", "x", backend.Stringify())
		} else {
			fmt.Printf(" [%s] %s", " ", backend.Stringify())
		}
	}
}

func NewStaticBS(backends []loadbalancer.Backend) *StaticBS {
	strategy := new(StaticBS)
	strategy.Init(backends)
	return strategy
}
