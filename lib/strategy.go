package lib

import (
	"crypto/sha256"
	"fmt"
	"math/big"
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

func (thbs *TraditionalHasingBS) hashFn(key string) *big.Int {
	h := sha256.New()

	h.Write([]byte(key))
	md := h.Sum(nil)

	keymodbackends := new(big.Int)

	mdBigInt := new(big.Int).SetBytes(md)

	totalBackends := len(thbs.Backends)
	return keymodbackends.Mod(mdBigInt, big.NewInt(int64(totalBackends)))
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
	backendIndex := tbhs.hashFn(req.reqId)

	return tbhs.Backends[backendIndex.Int64()]
}

func (tbhs *TraditionalHasingBS) RegisterBackend(backend *Backend) {
	tbhs.Backends = append(tbhs.Backends, backend)
}

func (thbs *TraditionalHasingBS) PrintTopology() {
	for _, backend := range thbs.Backends {
		fmt.Printf("	[%s] %s", " ", backend.stringfy())
	}
}
