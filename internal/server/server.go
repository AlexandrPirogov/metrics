package server

import (
	"context"
	"log"
	"memtracker/internal/config/server"
	"memtracker/internal/server/db"
	"memtracker/internal/server/http"
	"memtracker/internal/server/rpc"
)

type MetricServer interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

func NewMetricServer() MetricServer {
	log.Println("Running http server")
	return http.BuildHTPP()
}

func NewRPC() MetricServer {
	log.Println("Running rpc server")
	return &rpc.RPCServer{
		Storer: db.GetStorer(),
	}
}

func New() MetricServer {
	log.Println("rpc ", server.ServerCfg.RPC)
	if !server.ServerCfg.RPC {
		return NewMetricServer()
	}
	return NewRPC()
}
