package rpc

import (
	"context"
	fmt "fmt"
	"log"
	"memtracker/internal/config/agent"
	"memtracker/internal/crypt"
	"memtracker/internal/function"
	"memtracker/internal/metrics"

	grpc "google.golang.org/grpc"
)

type RPCClient struct {
	workers int
	stream  MetricHandler_UpdateGaugesClient
	conn    MetricHandlerClient
}

func New() RPCClient {
	conn, err := grpc.Dial(agent.ClientCfg.Address, grpc.WithInsecure())
	function.ErrFatalCheck("can't connect via RPC to server: ", err)

	clientConn := NewMetricHandlerClient(conn)
	stream, err := clientConn.UpdateGauges(context.Background())
	function.ErrFatalCheck("can't establish stream: ", err)

	return RPCClient{
		stream: stream,
		conn:   clientConn,
	}
}

func (r RPCClient) Send(metrics []metrics.Metricable) {
	for _, m := range metrics {
		if m.String() == "gauge" {
			gauges := buildMessageGauges("gauge", m.AsMap())
			for _, g := range gauges {
				r.AddInStream(&g, r.stream)
			}
		} else {

		}
	}
}

func (r RPCClient) Listen() {

}

func buildMessageGauges(typ string, gauges map[string]interface{}) []Gauge {
	key := agent.ClientCfg.Hash
	res := make([]Gauge, 0)
	for k, v := range gauges {
		val := v.(float64)
		toMarsal := Gauge{
			Id:    k,
			Type:  typ,
			Value: val,
			Hash:  crypt.Hash(fmt.Sprintf("%s:gauge:%f", k, val), key),
		}
		res = append(res, toMarsal)
	}
	return res
}

func (r RPCClient) AddInStream(g *Gauge, stream MetricHandler_UpdateGaugesClient) {
	if err := stream.Send(g); err != nil {
		log.Fatal(err)
	}
}
