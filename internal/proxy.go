package internal

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/maniSHarma7575/loadbalancer/internal/loadbalancer"
	"github.com/maniSHarma7575/loadbalancer/internal/strategy"
)

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

type Event struct {
	EventName string
	Data      interface{}
}

type LB struct {
	Backends []loadbalancer.Backend
	Events   chan Event
	Strategy loadbalancer.BalancingStrategy
}

type IncomingReq struct {
	SrcConn net.Conn
	ReqId   string
}

func (req *IncomingReq) GetReqID() string {
	return req.ReqId
}

func (req *IncomingReq) GetSrcConn() net.Conn {
	return req.SrcConn
}

var lb *LB

func InitLB() *LB {
	backends := []loadbalancer.Backend{
		&Backend{Host: "localhost", Port: 8085, IsHealthy: true},
		&Backend{Host: "localhost", Port: 8086, IsHealthy: true},
		&Backend{Host: "localhost", Port: 8087, IsHealthy: true},
	}

	lb = &LB{
		Backends: backends,
		Events:   make(chan Event),
		Strategy: strategy.NewConsistentHashingBS(backends),
	}
	return lb
}

func (lb *LB) Run() {
	listener, err := net.Listen("tcp", ":8082")

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	log.Println("LB listening on port 8082")

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

func (lb *LB) Proxy(req loadbalancer.IncomingReq) {
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
