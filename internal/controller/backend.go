package controller

import (
	"fmt"
)

type Backend struct {
	Port            int
	Host            string
	IsHealthy       bool
	NumberRequests  int
	HealthStatusUrl string
}

func (b *Backend) Stringify() string {
	return fmt.Sprintf("http://%s:%d", b.Host, b.Port)
}

func (b *Backend) IncrementRequestCounter() {
	b.NumberRequests++
}

func (b *Backend) GetHealthStatusUrl() string {
	return b.HealthStatusUrl
}

func (b *Backend) UpdateIsHealthy(status bool) {
	b.IsHealthy = status
}

func (b *Backend) IsBackendHealthy() bool {
	return b.IsHealthy
}
