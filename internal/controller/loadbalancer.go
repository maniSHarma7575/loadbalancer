package controller

import (
	"io"
	"log"
	"net"
	"slices"
	"sync"

	"github.com/google/uuid"
	loadbalancer "github.com/maniSHarma7575/loadbalancer/internal/balancer"
	"github.com/maniSHarma7575/loadbalancer/internal/strategy"
)

type LoadBalancer struct {
	Backends []loadbalancer.Backend
	Events   chan loadbalancer.Event
	Strategy loadbalancer.BalancingStrategy
	sync.RWMutex
}

var lb *LoadBalancer

func InitLB(configs map[string]interface{}) *LoadBalancer {
	backends := []loadbalancer.Backend{}

	for _, backendDetails := range configs["Backends"].([]map[string]interface{}) {
		backends = append(backends, &Backend{
			Host:            backendDetails["Host"].(string),
			Port:            backendDetails["Port"].(int),
			IsHealthy:       backendDetails["IsHealthy"].(bool),
			HealthStatusUrl: backendDetails["HealthStatusUrl"].(string),
		})
	}

	lb = &LoadBalancer{
		Backends: backends,
		Events:   make(chan loadbalancer.Event),
		Strategy: strategy.NewConsistentHashingBS(backends),
	}

	lb.ChangeStrategy(configs["Strategy"].(string))
	return lb
}

func (lb *LoadBalancer) Run() {
	healthCheckInterval := 3
	healthChecker := NewHealthChecker(lb.Backends)
	healthChecker.Attach(lb)
	healthChecker.Start(healthCheckInterval)

	listener, err := net.Listen("tcp", ":8082")

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	log.Println("LB listening on port 8082")

	go lb.RunEventLoop()

	for {
		connection, err := listener.Accept()

		if err != nil {
			log.Printf("unable to accept the connection: %s", err.Error())
		}

		go lb.Proxy(&IncomingReq{
			SrcConn: connection,
			ReqId:   uuid.NewString(),
		})
	}
}

func (lb *LoadBalancer) AddBackend(backend *Backend) {
	lb.Backends = append(lb.Backends, backend)
	lb.Strategy.RegisterBackend(backend)
}

func (lb *LoadBalancer) ChangeStrategy(stratgeyName string) {
	switch stratgeyName {
	case "round-robin":
		lb.Strategy = strategy.NewRoundRobinBS(lb.Backends)
	case "static":
		lb.Strategy = strategy.NewStaticBS(lb.Backends)
	case "traditional_hash":
		lb.Strategy = strategy.NewTraditionalHashBS(lb.Backends)
	case "consistent_hash":
		lb.Strategy = strategy.NewConsistentHashingBS(lb.Backends)
	default:
		lb.Strategy = strategy.NewConsistentHashingBS(lb.Backends)
	}
}

func (lb *LoadBalancer) RunEventLoop() {
	for {
		select {
		case event := <-lb.Events:
			switch event.GetEventName() {
			case EXIT:
				log.Println("Gracefully terminating")
				return
			case ADD_BACKEND:
				backend, ok := event.GetData().(Backend)
				if !ok {
					panic("Something wrong with you backend!")
				}
				lb.AddBackend(&backend)
			case CHANGE_STRATEGY:
				strategyName, ok := event.GetData().(string)
				if !ok {
					panic("Please give strategy name in string format")
				}
				lb.ChangeStrategy(strategyName)
			}
		}
	}
}

func (lb *LoadBalancer) Proxy(req loadbalancer.IncomingReq) {
	backend := lb.Strategy.GetNextBackend(req)

	log.Printf("in-req: %s out-req: %s", req.GetReqID(), backend.Stringify())

	backendConn, err := net.Dial("tcp", backend.Stringify())

	if err != nil {
		log.Printf("error connecting to backend: %s", err.Error())
		req.GetSrcConn().Write([]byte("backend not avaiable"))
		req.GetSrcConn().Close()
		panic(err)
	}

	backend.IncrementRequestCounter()

	go io.Copy(backendConn, req.GetSrcConn())
	go io.Copy(req.GetSrcConn(), backendConn)
}

func (lb *LoadBalancer) BackendUp(backend loadbalancer.Backend) {
	defer lb.Unlock()

	lb.Lock()

	idx := slices.IndexFunc(lb.Backends, func(b loadbalancer.Backend) bool { return backend == b })
	if idx != -1 {
		lb.Backends[idx].UpdateIsHealthy(true)
		lb.Strategy.RefreshBackend(lb.Backends[idx])
	}
	log.Printf("Server is up: %s", backend.Stringify())
}

func (lb *LoadBalancer) BackendDown(backend loadbalancer.Backend) {
	defer lb.Unlock()

	lb.Lock()

	idx := slices.IndexFunc(lb.Backends, func(b loadbalancer.Backend) bool { return backend == b })
	if idx != -1 {
		lb.Backends[idx].UpdateIsHealthy(false)
		lb.Strategy.RefreshBackend(lb.Backends[idx])
	}
	log.Printf("Server went down: %s", backend.Stringify())
}
