package server

import (
	"context"
	"memtracker/internal/config/server"
	"memtracker/internal/server/db"
	"memtracker/internal/server/http"
	"memtracker/internal/server/rpc"
)

type MetricServer interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

var isRPC map[bool]func() MetricServer = map[bool]func() MetricServer{
	true:  NewRPC,
	false: NewMetricServer,
}

func NewMetricServer() MetricServer {
	return http.BuildHTPP()
}

func NewRPC() MetricServer {
	return &rpc.RPCServer{
		Storer: db.GetStorer(),
	}
}

func New() MetricServer {
	server := server.ServerCfg.RPC
	return isRPC[server]()
}
