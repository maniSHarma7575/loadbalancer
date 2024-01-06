package strategy

import (
	"fmt"

	loadbalancer "github.com/maniSHarma7575/loadbalancer/internal/balancer"
)

type TraditionalHasingBS struct {
	Backends []loadbalancer.Backend
}

func NewTraditionalHashBS(backends []loadbalancer.Backend) *TraditionalHasingBS {
	strategy := new(TraditionalHasingBS)
	strategy.Init(backends)
	return strategy
}

func (tbhs *TraditionalHasingBS) Init(backends []loadbalancer.Backend) {
	tbhs.Backends = backends
}

func (tbhs *TraditionalHasingBS) GetNextBackend(req loadbalancer.IncomingReq) loadbalancer.Backend {
	backendIndex := hashFn(req.GetReqID()).Int64() % int64(len(tbhs.Backends))

	return tbhs.Backends[backendIndex]
}

func (tbhs *TraditionalHasingBS) RegisterBackend(backend loadbalancer.Backend) {
	tbhs.Backends = append(tbhs.Backends, backend)
}

func (thbs *TraditionalHasingBS) PrintTopology() {
	for _, backend := range thbs.Backends {
		fmt.Printf("	[%s] %s", " ", backend.Stringify())
	}
}
