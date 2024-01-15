package strategy

import (
	"fmt"

	loadbalancer "github.com/maniSHarma7575/loadbalancer/internal/balancer"
)

type RoundRobinBS struct {
	Index    int
	Backends []loadbalancer.Backend
}

func (rrbs *RoundRobinBS) Init(backends []loadbalancer.Backend) {
	rrbs.Index = 0
	rrbs.Backends = backends
}

func (rrbs *RoundRobinBS) GetNextBackend(loadbalancer.IncomingReq) loadbalancer.Backend {
	rrbs.Index = (rrbs.Index + 1) % len(rrbs.Backends)
	loopCounter := 0

	for loopCounter < len(rrbs.Backends) && !rrbs.Backends[rrbs.Index].IsBackendHealthy() {
		rrbs.Index = (rrbs.Index + 1) % len(rrbs.Backends)
		loopCounter += 1
	}
	return rrbs.Backends[rrbs.Index]
}

func (rrbs *RoundRobinBS) RegisterBackend(backend loadbalancer.Backend) {
	rrbs.Backends = append(rrbs.Backends, backend)
}

func (rrbs *RoundRobinBS) PrintTopology() {
	for index, backend := range rrbs.Backends {
		fmt.Printf("		[%d] %s", index, backend.Stringify())
	}
}

func NewRoundRobinBS(backends []loadbalancer.Backend) *RoundRobinBS {
	strategy := new(RoundRobinBS)
	strategy.Init(backends)
	return strategy
}
