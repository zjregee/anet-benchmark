package main

import (
	"context"

	"anet-benchmark/runner"
	"anet-benchmark/netpoll/codec"
	"github.com/cloudwego/netpoll"
)

func NewRPCServer(network, address string) runner.Server {
	return &rpcServer{
		network: network,
		address: address,
	}
}

var _ runner.Server = &rpcServer{}

type rpcServer struct{
	network string
	address string
}

func (s *rpcServer) Run() error {
	listener, err := netpoll.CreateListener(s.network, s.address)
	if err != nil {
		panic("can't failed here")
	}
	eventLoop, err := netpoll.NewEventLoop(s.handler)
	if err != nil {
		panic("can't failed here")
	}
	return eventLoop.Serve(listener)
}

func (s *rpcServer) handler(ctx context.Context, conn netpoll.Connection) error {
	reader, writer := conn.Reader(), conn.Writer()
	req := &runner.Message{}
	err := codec.Decode(reader, req)
	if err != nil {
		return err
	}
	resp := &runner.Message{ Message: req.Message }
	err = codec.Encode(writer, resp)
	return err
}
