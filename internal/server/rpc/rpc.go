package rpc

import (
	"context"
	"errors"
	fmt "fmt"
	"io"
	"log"
	"memtracker/internal/config/server"
	"memtracker/internal/crypt"
	"memtracker/internal/function"
	"memtracker/internal/kernel"
	"memtracker/internal/kernel/tuples"
	"memtracker/internal/server/db"
	"memtracker/internal/server/verifier"
	"net"

	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
)

type RPCServer struct {
	Storer db.DB
}

func newRPC() *grpc.Server {
	if server.ServerCfg.CryptoKey != server.DefaultCryptoKey {
		creds, err := server.LoadRPCTLSCredentials(server.ServerCfg.CryptoKey)
		function.ErrFatalCheck("", err)
		log.Println("Turning TLS...")
		return grpc.NewServer(
			grpc.Creds(creds),
		)
	}

	return grpc.NewServer()
}

func (g *RPCServer) ListenAndServe() error {
	listener, err := net.Listen("tcp", server.ServerCfg.Address)
	function.ErrFatalCheck("can't start rpc server", err)
	s := newRPC()

	RegisterMetricHandlerServer(s, g)
	return s.Serve(listener)
}

func (g *RPCServer) Shutdown(ctx context.Context) error {
	return g.Shutdown(ctx)
}

func (g *RPCServer) UpdateGauges(stream MetricHandler_UpdateGaugesServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())

	if !verifier.IsContainerSubnetGRPC(md) {
		log.Println("subnet not allowed")
		return errors.New("asd")
	}
	var err error
	for ok {
		err = g.recieveGauges(stream)
	}
	return err
}

func (g *RPCServer) UpdateCounters(stream MetricHandler_UpdateCountersServer) error {
	for {
		g.recieveCounters(stream)
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

func (g *RPCServer) recieveCounters(stream MetricHandler_UpdateCountersServer) error {
	counter, err := stream.Recv()
	if err == io.EOF {
		return stream.SendAndClose(
			&wrappers.StringValue{Value: "Processed counters "})
	}

	if err != nil {
		return stream.SendAndClose(
			&wrappers.StringValue{Value: "Counters processed " + fmt.Sprint(err)})
	}

	tupl := counter.ToTuple()
	log.Println(tupl)
	kernel.Write(g.Storer.Storage, tuples.TupleList{}.Add(tupl))
	return nil
}

func (g Gauge) ToTuple() tuples.Tupler {
	res := tuples.NewTuple()
	res.SetField("name", g.Id)
	res.SetField("type", g.Type)
	res.SetField("hash", crypt.Hash(fmt.Sprintf("%s:gauge:%f", g.Id, g.Value), server.ServerCfg.Hash))
	res.SetField("value", &g.Value)
	return res
}

func (c Counter) ToTuple() tuples.Tupler {
	res := tuples.NewTuple()
	res.SetField("name", c.Id)
	res.SetField("type", c.Type)
	res.SetField("hash", crypt.Hash(fmt.Sprintf("%s:counter:%d", c.Id, c.Delta), server.ServerCfg.Hash))
	res.SetField("value", &c.Delta)
	return res
}
