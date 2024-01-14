package controller

import (
	"net/http"
	"sync"
	"time"

	"github.com/maniSHarma7575/loadbalancer/internal/balancer"
)

type HealthChecker struct {
	Backends   []balancer.Backend
	Observers  []balancer.Observer
	Status     map[balancer.Backend]bool
	HttpClient *http.Client
	sync.RWMutex
}

func (h *HealthChecker) Start(interval int) {
	go func() {
		for {
			h.Check()
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}()
}

func (h *HealthChecker) Attach(observer balancer.Observer) {
	h.Observers = append(h.Observers, observer)
}

func (h *HealthChecker) Check() {
	for _, backend := range h.Backends {
		go h.notify(backend)
	}
}

func (h *HealthChecker) notify(backend balancer.Backend) {
	ishealthy := h.isHealthy(backend)
	isStatusChaged := h.statusChanged(backend, ishealthy)

	if isStatusChaged {
		if ishealthy {
			h.notifyBackendUp(backend)
		} else {
			h.notifyBackendDown(backend)
		}
	}

	h.changeStatus(backend, ishealthy)
}

func (h *HealthChecker) notifyBackendUp(backend balancer.Backend) {
	for _, observer := range h.Observers {
		observer.BackendUp(backend)
	}
}

func (h *HealthChecker) notifyBackendDown(backend balancer.Backend) {
	for _, observer := range h.Observers {
		observer.BackendDown(backend)
	}
}

func (h *HealthChecker) changeStatus(backend balancer.Backend, isHealthy bool) {
	defer h.Unlock()

	h.Lock()
	h.Status[backend] = isHealthy
}

func (h *HealthChecker) statusChanged(backend balancer.Backend, isHealthy bool) bool {
	defer h.RUnlock()

	h.RLock()

	status, ok := h.Status[backend]
	if !ok {
		return true
	}

	return status != isHealthy
}

func (h *HealthChecker) isHealthy(backend balancer.Backend) bool {
	res, err := h.HttpClient.Get(backend.GetHealthStatusUrl())

	if err == nil && res.StatusCode == http.StatusOK {
		return true
	}

	return false
}

func NewHealthChecker(backends []balancer.Backend) balancer.HealthCheckerInterface {
	return &HealthChecker{
		Backends:   backends,
		HttpClient: http.DefaultClient,
		Status:     make(map[balancer.Backend]bool),
	}
}
