package strategy

import (
	"fmt"
	"slices"
	"sort"
	"sync"

	loadbalancer "github.com/maniSHarma7575/loadbalancer/internal/balancer"
)

type ConsistentHashingBS struct {
	Backends []loadbalancer.Backend
	Slots    []int
	sync.RWMutex
}

func NewConsistentHashingBS(backends []loadbalancer.Backend) *ConsistentHashingBS {
	strategy := new(ConsistentHashingBS)
	strategy.Init(backends)
	return strategy
}

func (chbs *ConsistentHashingBS) Init(backends []loadbalancer.Backend) {
	chbs.Slots = []int{}
	chbs.Backends = []loadbalancer.Backend{}

	for _, backend := range backends {
		backendHash := hashFn(backend.Stringify()).Int64()

		if len(chbs.Slots) == 0 {
			chbs.Slots = append(chbs.Slots, int(backendHash))
			chbs.Backends = append(chbs.Backends, backend)
			continue
		}

		index := sort.Search(len(chbs.Slots), func(i int) bool {
			return chbs.Slots[i] >= int(backendHash)
		})

		if index == len(chbs.Slots) {
			chbs.Slots = append(chbs.Slots, int(backendHash))
		} else {
			chbs.Slots = append(chbs.Slots[:index+1], chbs.Slots[index:]...)
			chbs.Slots[index] = int(backendHash)
		}

		if index == len(chbs.Backends) {
			chbs.Backends = append(chbs.Backends, backend)
		} else {
			chbs.Backends = append(chbs.Backends[:index+1], chbs.Backends[index:]...)
			chbs.Backends[index] = backend
		}
	}
}

func (chbs *ConsistentHashingBS) RegisterBackend(backend loadbalancer.Backend) {
	backendHash := hashFn(backend.Stringify()).Int64()

	index := sort.Search(len(chbs.Slots), func(i int) bool {
		return chbs.Slots[i] >= int(backendHash)
	})

	if index == len(chbs.Slots) {
		chbs.Slots = append(chbs.Slots, int(backendHash))
	} else {
		chbs.Slots = append(chbs.Slots[:index+1], chbs.Slots[index:]...)
		chbs.Slots[index] = int(backendHash)
	}

	if index == len(chbs.Backends) {
		chbs.Backends = append(chbs.Backends, backend)
	} else {
		chbs.Backends = append(chbs.Backends[:index+1], chbs.Backends[index:]...)
		chbs.Backends[index] = backend
	}
}

func (chbs *ConsistentHashingBS) RefreshBackend(backend loadbalancer.Backend) {
	defer chbs.Unlock()

	chbs.Lock()
	if !backend.IsBackendHealthy() {
		idx := FindBackendIndex(chbs.Backends, backend)

		if idx != -1 {
			chbs.Backends = append(chbs.Backends[:idx], chbs.Backends[idx+1:]...)
			backendHash := hashFn(backend.Stringify()).Int64()

			slotsIndex := slices.IndexFunc(chbs.Slots, func(b int) bool {
				return int(backendHash) == b
			})

			if slotsIndex != -1 {
				chbs.Slots = append(chbs.Slots[:slotsIndex], chbs.Slots[slotsIndex+1:]...)
			}
		}
	} else {
		idx := FindBackendIndex(chbs.Backends, backend)

		if idx == -1 {
			chbs.RegisterBackend(backend)
		}
	}
}

func (chbs *ConsistentHashingBS) GetNextBackend(req loadbalancer.IncomingReq) loadbalancer.Backend {
	hash := hashFn(req.GetReqID()).Int64()

	index := sort.Search(len(chbs.Slots), func(i int) bool {
		return chbs.Slots[i] >= int(hash)
	})

	return chbs.Backends[index%len(chbs.Backends)]
}

func (chbs *ConsistentHashingBS) PrintTopology() {
	index, i := 0, 0
	for i = 0; i < 19; i++ {
		if index < len(chbs.Slots) && chbs.Slots[index] == i {
			fmt.Printf("		[%2d] %s", i, chbs.Backends[index].Stringify())
			index++
		} else {
			fmt.Printf("		[%2d] -", i)
		}
	}
}
