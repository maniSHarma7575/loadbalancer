package lib

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/google/uuid"
)

type Backend struct {
	Port           int
	Host           string
	IsHealthy      bool
	NumberRequests int
}

func (b *Backend) stringfy() string {
	return fmt.Sprintf("%s:%d", b.Host, b.Port)
}

type Event struct {
	EventName string
	Data      interface{}
}

type LB struct {
	Backends []*Backend
	Events   chan Event
	Strategy BalancingStrategy
}

type IncomingReq struct {
	srcConn net.Conn
	reqId   string
}

var lb *LB

func InitLB() *LB {
	backends := []*Backend{
		&Backend{Host: "localhost", Port: 8085, IsHealthy: true},
		&Backend{Host: "localhost", Port: 8086, IsHealthy: true},
		&Backend{Host: "localhost", Port: 8087, IsHealthy: true},
	}

	lb = &LB{
		Backends: backends,
		Events:   make(chan Event),
		Strategy: NewRoundRobinBS(backends),
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

		go lb.Proxy(IncomingReq{
			srcConn: connection,
			reqId:   uuid.NewString(),
		})
	}
}

func (lb *LB) Proxy(req IncomingReq) {
	backend := lb.Strategy.GetNextBackend(req)

	log.Printf("in-req: %s out-req: %s", req.reqId, backend.stringfy())

	backendConn, err := net.Dial("tcp", backend.stringfy())

	if err != nil {
		log.Printf("error connecting to backend: %s", err.Error())
		req.srcConn.Write([]byte("backend not avaiable"))
		req.srcConn.Close()
		panic(err)
	}

	backend.NumberRequests++

	go io.Copy(backendConn, req.srcConn)
	go io.Copy(req.srcConn, backendConn)
}
