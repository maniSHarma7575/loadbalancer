package controller

import (
	"io"
	"log"
	"net"
	"slices"
	"strconv"
	"sync"

	"github.com/google/uuid"
	loadbalancer "github.com/maniSHarma7575/loadbalancer/internal/balancer"
	"github.com/maniSHarma7575/loadbalancer/internal/config"
	"github.com/maniSHarma7575/loadbalancer/internal/strategy"
)

type LoadBalancer struct {
	Backends []loadbalancer.Backend
	Events   chan loadbalancer.Event
	Strategy loadbalancer.BalancingStrategy
	Config   config.Config
	sync.RWMutex
}

var lb *LoadBalancer

func InitLB(configs map[string]interface{}) *LoadBalancer {
	var cfg config.Config
	configPaths := config.ConfigPaths()
	for _, configPath := range configPaths {
		config, err := config.Load(configPath)
		if err == nil {
			cfg = *config
			break
		}
	}

	backends := []loadbalancer.Backend{}

	for _, server := range *cfg.Servers {
		backends = append(backends, &Backend{
			Host:            server.Host,
			Port:            server.Port,
			IsHealthy:       false,
			HealthStatusUrl: "http://" + server.Host + ":" + strconv.Itoa(server.Port) + server.HealthPath,
		})
	}

	lb = &LoadBalancer{
		Backends: backends,
		Events:   make(chan loadbalancer.Event),
		Strategy: strategy.NewConsistentHashingBS(backends),
		Config:   cfg,
	}

	lb.ChangeStrategy(cfg.LoadBalanceStrategy)
	return lb
}

func (lb *LoadBalancer) Run() {
	healthCheckInterval := lb.Config.HealthCheckIntervalSeconds
	healthChecker := NewHealthChecker(lb.Backends)
	healthChecker.Attach(lb)
	healthChecker.Start(healthCheckInterval)

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(lb.Config.Port))

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
