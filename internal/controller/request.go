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

func (req *IncomingReq) GetHeadersAsMap() map[string]string {
	converted := make(map[string]string)
	for key, values := range req.Request.Header {
		if len(values) > 0 {
			converted[key] = values[0]
		}
	}
	return converted
}
