package controller

import "net"

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
