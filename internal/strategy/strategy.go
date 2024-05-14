package strategy

import (
	"crypto/md5"
	"math/big"
	"slices"

	"github.com/maniSHarma7575/loadbalancer/internal/balancer"
)

func hashFn(key string) *big.Int {
	h := md5.New()

	h.Write([]byte(key))
	md := h.Sum(nil)

	keymodbackends := new(big.Int)

	mdBigInt := new(big.Int).SetBytes(md)

	return keymodbackends.Mod(mdBigInt, big.NewInt(int64(19)))
}

func HelathyBackends(backends []balancer.Backend) []balancer.Backend {
	healthyBackends := make([]balancer.Backend, 0)
	for _, backend := range backends {
		if backend.IsBackendHealthy() {
			healthyBackends = append(healthyBackends, backend)
		}
	}
	return healthyBackends
}

func FindBackendIndex(backends []balancer.Backend, backend balancer.Backend) int {
	idx := slices.IndexFunc(backends, func(b balancer.Backend) bool { return backend == b })
	return idx
}

var (
	RoundRobinStrategy         = "round-robin"
	StaticStrategy             = "static"
	TraditionalHashingStrategy = "traditional_hash"
	ConsistentHashingStrategy  = "consistent_hash"
	StickySessionStrategy      = "sticky_session"
	LeastConnectionStrategy    = "least_connections"
)
