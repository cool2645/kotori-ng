package main

import (
	"github.com/cool2645/kotori-ng/handler"
)

const (
	BaseApiVer = "/v1"
	Base       = BaseApi + BaseApiVer
)

var (
	v1Api      = api.PathPrefix(BaseApiVer).Subrouter()
)

func RegRouter() {
	// Ping
	v1Api.Methods("GET").Path("").HandlerFunc(handler.Pong)
}
