package strategy

import (
	"fmt"
	"math"
	"sync"

	loadbalancer "github.com/maniSHarma7575/loadbalancer/internal/balancer"
)

type LeastConnectionBS struct {
	Backends         []loadbalancer.Backend
	ConnectionsCount map[loadbalancer.Backend]int64
	sync.RWMutex
}

func NewLeastConnectionBS(backends []loadbalancer.Backend) *LeastConnectionBS {
	strategy := new(LeastConnectionBS)
	strategy.Init(backends)
	return strategy
}

func (lcbs *LeastConnectionBS) Init(backends []loadbalancer.Backend) {
	lcbs.Backends = backends
	lcbs.ConnectionsCount = make(map[loadbalancer.Backend]int64)
	for _, backend := range backends {
		lcbs.ConnectionsCount[backend] = 0
	}
}

func (lcbs *LeastConnectionBS) RegisterBackend(backend loadbalancer.Backend) {
	lcbs.Backends = append(lcbs.Backends, backend)
	lcbs.ConnectionsCount[backend] = 0
}

func (lcbs *LeastConnectionBS) PrintTopology() {
	for _, backend := range lcbs.Backends {
		fmt.Printf("	[%s] %s", " ", backend.Stringify())
	}
}

func (lcbs *LeastConnectionBS) GetNextBackend(loadbalancer.IncomingReq) loadbalancer.Backend {
	defer lcbs.Unlock()

	lcbs.Lock()
	minCount := int64(math.MaxInt64)

	var selectedBackend loadbalancer.Backend

	for backend, count := range lcbs.ConnectionsCount {
		if count < minCount {
			minCount = count
			selectedBackend = backend
		}
	}

	lcbs.increaseConnectionCount(selectedBackend)

	return selectedBackend
}

func (lcbs *LeastConnectionBS) increaseConnectionCount(backend loadbalancer.Backend) {
	if _, exists := lcbs.ConnectionsCount[backend]; !exists {
		lcbs.ConnectionsCount[backend] = 0
	}

	lcbs.ConnectionsCount[backend]++
}

func (lcbs *LeastConnectionBS) RefreshBackend(backend loadbalancer.Backend) {
	defer lcbs.Unlock()

	lcbs.Lock()
	idx := FindBackendIndex(lcbs.Backends, backend)

	if backend.IsBackendHealthy() && idx == -1 {
		lcbs.Backends = append(lcbs.Backends, backend)
		lcbs.ConnectionsCount[backend] = 0
	} else if idx != -1 && !backend.IsBackendHealthy() {
		lcbs.Backends = append(lcbs.Backends[:idx], lcbs.Backends[idx+1:]...)
		delete(lcbs.ConnectionsCount, backend)
	}
}
