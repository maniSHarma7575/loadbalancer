package lib

import (
	"crypto/md5"
	"fmt"
	"math/big"
	"sort"
)

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

type ConsistentHashingBS struct {
	Backends []*Backend
	Slots    []int
}

type TraditionalHasingBS struct {
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
		fmt.Printf("		[%d] %s", index, backend.stringfy())
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
			fmt.Printf("	[%s] %s", "x", backend.stringfy())
		} else {
			fmt.Printf(" [%s] %s", " ", backend.stringfy())
		}
	}
}

func NewStaticBS(backends []*Backend) *StaticBS {
	strategy := new(StaticBS)
	strategy.Init(backends)
	return strategy
}

func hashFn(key string) *big.Int {
	h := md5.New()

	h.Write([]byte(key))
	md := h.Sum(nil)

	keymodbackends := new(big.Int)

	mdBigInt := new(big.Int).SetBytes(md)

	return keymodbackends.Mod(mdBigInt, big.NewInt(int64(19)))
}

func NewTraditionalHashBS(backends []*Backend) *TraditionalHasingBS {
	strategy := new(TraditionalHasingBS)
	strategy.Init(backends)
	return strategy
}

func (tbhs *TraditionalHasingBS) Init(backends []*Backend) {
	tbhs.Backends = backends
}

func (tbhs *TraditionalHasingBS) GetNextBackend(req IncomingReq) *Backend {
	backendIndex := hashFn(req.reqId).Int64() % int64(len(tbhs.Backends))

	return tbhs.Backends[backendIndex]
}

func (tbhs *TraditionalHasingBS) RegisterBackend(backend *Backend) {
	tbhs.Backends = append(tbhs.Backends, backend)
}

func (thbs *TraditionalHasingBS) PrintTopology() {
	for _, backend := range thbs.Backends {
		fmt.Printf("	[%s] %s", " ", backend.stringfy())
	}
}

func NewConsistentHashingBS(backends []*Backend) *ConsistentHashingBS {
	stratgey := new(ConsistentHashingBS)
	stratgey.Init(backends)
	return stratgey
}

func (chbs *ConsistentHashingBS) Init(backends []*Backend) {
	chbs.Slots = []int{}
	chbs.Backends = []*Backend{}

	for _, backend := range backends {
		backendHash := hashFn(backend.stringfy()).Int64()

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

func (chbs *ConsistentHashingBS) RegisterBackend(backend *Backend) {
	backendHash := hashFn(backend.stringfy()).Int64()

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

func (chbs *ConsistentHashingBS) GetNextBackend(req IncomingReq) *Backend {
	hash := hashFn(req.reqId).Int64()

	index := sort.Search(len(chbs.Slots), func(i int) bool {
		return chbs.Slots[i] >= int(hash)
	})

	return chbs.Backends[index%len(chbs.Backends)]
}

func (chbs *ConsistentHashingBS) PrintTopology() {
	index, i := 0, 0
	for i = 0; i < 19; i++ {
		if index < len(chbs.Slots) && chbs.Slots[index] == i {
			fmt.Printf("		[%2d] %s", i, chbs.Backends[index].stringfy())
			index++
		} else {
			fmt.Printf("		[%2d] -", i)
		}
	}
}
