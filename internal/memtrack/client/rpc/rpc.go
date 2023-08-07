package rpc

import (
	context "context"
	fmt "fmt"
	"log"
	"memtracker/internal/config/agent"
	"memtracker/internal/crypt"
	"memtracker/internal/function"
	"memtracker/internal/metrics"

	"github.com/net-byte/go-gateway"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
)

var buildAndSendByType map[string]func(RPCClient, metrics.Metricable) = map[string]func(RPCClient, metrics.Metricable){
	"gauge":   buildAndSendGauges,
	"counter": buildAndSendCounters,
}

type RPCClient struct {
	workers       int
	gaugeStream   MetricHandler_UpdateGaugesClient
	counterStream MetricHandler_UpdateCountersClient
	conn          MetricHandlerClient
}

func New() RPCClient {
	conn, err := grpc.Dial(agent.ClientCfg.Address, grpc.WithInsecure())
	function.ErrFatalCheck("can't connect via RPC to server: ", err)

	clientConn := NewMetricHandlerClient(conn)

	ip, _ := gateway.DiscoverGatewayIPv4()
	md := metadata.Pairs("X-Real-IP", ip.String())
	gctx := metadata.NewOutgoingContext(context.Background(), md)
	gstream, err := clientConn.UpdateGauges(gctx)

	cstream, err := clientConn.UpdateCounters(gctx)

	function.ErrFatalCheck("can't establish stream: ", err)
	grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name))
	return RPCClient{
		gaugeStream:   gstream,
		counterStream: cstream,
		conn:          clientConn,
	}
}

func (r RPCClient) Send(metrics []metrics.Metricable) {
	for _, m := range metrics {
		buildAndSendByType[m.String()](r, m)
	}
}

func buildAndSendGauges(r RPCClient, m metrics.Metricable) {
	gauges := buildMessageGauges("gauge", m.AsMap())
	for _, g := range gauges {
		r.AddInStreamGauge(&g, r.gaugeStream)
	}
}

func buildAndSendCounters(r RPCClient, m metrics.Metricable) {
	counters := buildMessageCounters("counter", m.AsMap())
	for _, c := range counters {
		r.AddInStreamCounter(&c, r.counterStream)
	}
}

func buildMessageCounters(typ string, counters map[string]interface{}) []Counter {
	key := agent.ClientCfg.Hash
	res := make([]Counter, 0)
	for k, v := range counters {
		val := v.(float64)
		del := int64(val)
		toMarsal := Counter{
			Id:    k,
			Type:  typ,
			Delta: del,
			Hash:  crypt.Hash(fmt.Sprintf("%s:counter:%d", k, del), key),
		}
		res = append(res, toMarsal)
	}
	return res
}

func (r RPCClient) AddInStreamCounter(c *Counter, stream MetricHandler_UpdateCountersClient) {
	if err := stream.Send(c); err != nil {
		log.Fatal(err)
	}
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

func (r RPCClient) AddInStreamGauge(g *Gauge, stream MetricHandler_UpdateGaugesClient) {
	if err := stream.Send(g); err != nil {
		log.Fatal(err)
	}
}
