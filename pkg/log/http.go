package log

import (
	"net/http"
)

const (
	HTTPKey = key("http")
)

type HTTP struct {
	Latency  float64   `json:"latency"`
	Error    string    `json:"error"`
	Request  *Request  `json:"request"`
	Response *Response `json:"response"`
}

type Request struct {
	Host   string      `json:"host"`
	Route  string      `json:"route"`
	Header http.Header `json:"header"`
	Param  string      `json:"param"`
}

type Response struct {
	Header   http.Header `json:"header"`
	Status   int         `json:"status"`
	Body     interface{} `json:"body"`
	RemoteIP string      `json:"remote_ip"`
}
