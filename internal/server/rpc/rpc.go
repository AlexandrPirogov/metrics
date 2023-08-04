package rpc

import (
	"context"
	fmt "fmt"
	"io"
	"log"
	"memtracker/internal/config/server"
	"memtracker/internal/crypt"
	"memtracker/internal/function"
	"memtracker/internal/kernel"
	"memtracker/internal/kernel/tuples"
	"memtracker/internal/server/db"
	"net"

	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
)

type RPCServer struct {
	Storer db.DB
}

func (g *RPCServer) ListenAndServe() error {
	listener, err := net.Listen("tcp", server.ServerCfg.Address)
	function.ErrFatalCheck("can't start rpc server", err)

	s := grpc.NewServer()
	RegisterMetricHandlerServer(s, g)
	return s.Serve(listener)
}

func (g *RPCServer) Shutdown(ctx context.Context) error {
	return g.Shutdown(ctx)
}

func (g *RPCServer) UpdateGauges(stream MetricHandler_UpdateGaugesServer) error {
	for {
		g.recieveGauges(stream)
	}
}

func (g *RPCServer) recieveGauges(stream MetricHandler_UpdateGaugesServer) error {
	gauge, err := stream.Recv()
	if err == io.EOF {
		return stream.SendAndClose(
			&wrappers.StringValue{Value: "Processed gauges "})
	}

	if err != nil {
		return stream.SendAndClose(
			&wrappers.StringValue{Value: "Orders processed " + fmt.Sprint(err)})
	}

	tupl := gauge.ToTuple()
	log.Println(tupl)
	kernel.Write(g.Storer.Storage, tuples.TupleList{}.Add(tupl))
	return nil
}

func (g Gauge) ToTuple() tuples.Tupler {
	res := tuples.NewTuple()
	res.SetField("name", g.Id)
	res.SetField("type", g.Type)
	res.SetField("hash", crypt.Hash(fmt.Sprintf("%s:gauge:%f", g.Id, g.Value), server.ServerCfg.Hash))
	res.SetField("delta", g.Value)
	return res
}
