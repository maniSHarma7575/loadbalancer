package strategy

import (
	"fmt"
	"sync"

	loadbalancer "github.com/maniSHarma7575/loadbalancer/internal/balancer"
)

type TraditionalHasingBS struct {
	Backends []loadbalancer.Backend
	sync.RWMutex
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

func (tbhs *TraditionalHasingBS) RefreshBackend(backend loadbalancer.Backend) {
	defer tbhs.Unlock()

	tbhs.Lock()
	idx := FindBackendIndex(tbhs.Backends, backend)

	if backend.IsBackendHealthy() && idx == -1 {
		tbhs.Backends = append(tbhs.Backends, backend)
	} else if idx != -1 && !backend.IsBackendHealthy() {
		tbhs.Backends = append(tbhs.Backends[:idx], tbhs.Backends[idx+1:]...)
	}
}

func (tbhs *TraditionalHasingBS) RegisterBackend(backend loadbalancer.Backend) {
	tbhs.Backends = append(tbhs.Backends, backend)
}

func (thbs *TraditionalHasingBS) PrintTopology() {
	for _, backend := range thbs.Backends {
		fmt.Printf("	[%s] %s", " ", backend.Stringify())
	}
}
