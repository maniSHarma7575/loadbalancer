package controller

import "net/http"

type IncomingReq struct {
	ReqId   string
	Request *http.Request
}

func (req *IncomingReq) GetReqID() string {
	return req.ReqId
}

func (req *IncomingReq) GetHttpRequest() *http.Request {
	return req.Request
}
