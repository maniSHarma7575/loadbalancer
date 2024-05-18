package controller

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	loadbalancer "github.com/maniSHarma7575/loadbalancer/internal/balancer"
	"github.com/maniSHarma7575/loadbalancer/internal/config"
	"github.com/maniSHarma7575/loadbalancer/internal/strategy"
	"github.com/maniSHarma7575/loadbalancer/internal/utils"
)

type LoadBalancer struct {
	Backends []loadbalancer.Backend
	Events   chan loadbalancer.Event
	Strategy loadbalancer.BalancingStrategy
	Config   config.Config
	sync.RWMutex
}

var lb *LoadBalancer

func InitLB() *LoadBalancer {
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

func (lb *LoadBalancer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	req := &IncomingReq{
		ReqId:   uuid.NewString(),
		Request: request,
	}

	backend := lb.Strategy.GetNextBackend(req)

	t := time.Now()
	log.Printf("in-req: %s out-req: %s", req.GetReqID(), backend.Stringify())

	parsedUrl, _ := url.Parse(backend.Stringify())
	proxy := httputil.NewSingleHostReverseProxy(parsedUrl)
	cookie, _ := request.Cookie(lb.Config.StickySession.CookieKey)
	http.SetCookie(writer, cookie)
	proxy.ServeHTTP(writer, request)
	backend.IncrementRequestCounter()
	log.Printf("request served in server: %s time: %s", backend.Stringify(), time.Since(t))
}

func (lb *LoadBalancer) Run() {
	healthCheckInterval := lb.Config.HealthCheckIntervalSeconds
	healthChecker := NewHealthChecker(lb.Backends)
	healthChecker.Attach(lb)
	healthChecker.Start(healthCheckInterval)

	go lb.RunEventLoop()

	var httpHandler http.Handler = lb
	port := strconv.Itoa(lb.Config.Port)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: httpHandler,
	}

	log.Printf("LB listening on port :%s\n", port)
	tlsEnabled := lb.Config.Tls.Enabled

	if tlsEnabled {
		certFile := lb.Config.Tls.CertFile
		keyFile := lb.Config.Tls.KeyFile

		if !utils.IsFileExists(certFile) {
			panic(fmt.Errorf("TLS cert filepath doesn't exist %v", certFile))
		}

		if !utils.IsFileExists((keyFile)) {
			panic(fmt.Errorf("TLS key filepath doesn't exist %v", keyFile))
		}

		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}
}

func (lb *LoadBalancer) AddBackend(backend *Backend) {
	lb.Backends = append(lb.Backends, backend)
	lb.Strategy.RegisterBackend(backend)
}

func (lb *LoadBalancer) ChangeStrategy(stratgeyName string) {
	switch stratgeyName {
	case strategy.RoundRobinStrategy:
		lb.Strategy = strategy.NewRoundRobinBS(lb.Backends)
	case strategy.StaticStrategy:
		lb.Strategy = strategy.NewStaticBS(lb.Backends)
	case strategy.TraditionalHashingStrategy:
		lb.Strategy = strategy.NewTraditionalHashBS(lb.Backends)
	case strategy.ConsistentHashingStrategy:
		lb.Strategy = strategy.NewConsistentHashingBS(lb.Backends)
	case strategy.StickySessionStrategy:
		lb.Strategy = strategy.NewStickySessionBS(
			lb.Backends,
			lb.Config.StickySession.CookieKey,
			lb.Config.StickySession.TTLSeconds,
		)
	case strategy.LeastConnectionStrategy:
		lb.Strategy = strategy.NewLeastConnectionBS(lb.Backends)
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
