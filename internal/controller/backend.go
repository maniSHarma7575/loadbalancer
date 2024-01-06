package controller

import "fmt"

type Backend struct {
	Port           int
	Host           string
	IsHealthy      bool
	NumberRequests int
}

func (b *Backend) Stringify() string {
	return fmt.Sprintf("%s:%d", b.Host, b.Port)
}

func (b *Backend) IncrementRequestCounter() {
	b.NumberRequests++
}
