package strategy

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	loadbalancer "github.com/maniSHarma7575/loadbalancer/internal/balancer"
)

var (
	defaultCookieName = "lbcookie"
	defaultTTLSeconds = 120
)

type StickySessionBS struct {
	CookieName string
	TTLSeconds int
	Backends   []loadbalancer.Backend
	sync.RWMutex
}

func (ssbs *StickySessionBS) Init(backends []loadbalancer.Backend) {
	ssbs.Backends = backends
}

func (ssbs *StickySessionBS) GetNextBackend(request loadbalancer.IncomingReq) loadbalancer.Backend {
	httpRequest := request.GetHttpRequest()
	cookie, err := httpRequest.Cookie(ssbs.CookieName)
	if err != nil || cookie.Value == "" {
		cookieValue := uuid.NewString()
		httpRequest.AddCookie(&http.Cookie{
			Name:     ssbs.CookieName,
			Value:    cookieValue,
			Expires:  time.Now().Add(time.Second * time.Duration(ssbs.TTLSeconds)),
			HttpOnly: true,
		})
	}

	cookie, _ = httpRequest.Cookie(ssbs.CookieName)
	backendIndex := hashFn(cookie.Value).Int64() % int64(len(ssbs.Backends))
	return ssbs.Backends[backendIndex]
}

func (ssbs *StickySessionBS) RefreshBackend(backend loadbalancer.Backend) {
	defer ssbs.Unlock()

	ssbs.Lock()
	idx := FindBackendIndex(ssbs.Backends, backend)

	if backend.IsBackendHealthy() && idx == -1 {
		ssbs.Backends = append(ssbs.Backends, backend)
	} else if idx != -1 && !backend.IsBackendHealthy() {
		ssbs.Backends = append(ssbs.Backends[:idx], ssbs.Backends[idx+1:]...)
	}
}

func (ssbs *StickySessionBS) RegisterBackend(backend loadbalancer.Backend) {
	ssbs.Backends = append(ssbs.Backends, backend)
}

func (ssbs *StickySessionBS) PrintTopology() {
	for _, backend := range ssbs.Backends {
		fmt.Printf("	[%s] %s", " ", backend.Stringify())
	}
}

func NewStickySessionBS(backends []loadbalancer.Backend, cookieName string, TTLSeconds int) *StickySessionBS {
	strategy := new(StickySessionBS)
	strategy.Init(backends)
	if cookieName == "" {
		cookieName = defaultCookieName
	}

	if TTLSeconds <= 0 {
		TTLSeconds = defaultTTLSeconds
	}

	strategy.CookieName = cookieName
	strategy.TTLSeconds = TTLSeconds
	return strategy
}
