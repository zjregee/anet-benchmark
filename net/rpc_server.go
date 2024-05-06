package main

import (
	"fmt"
	"net"
	"time"
	"strings"

	"anet-benchmark/runner"
	"anet-benchmark/net/codec"
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
	listener, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "closed") {
				return err
			}
			fmt.Println("an error occurred while accepting, wait 10ms to retry")
			time.Sleep(10 * time.Millisecond)
			continue
		}
		go func(_conn net.Conn) {
			conner := codec.NewConner(_conn)
			defer codec.PutConner(conner)
			var err error
			for err == nil {
				err = s.handler(conner)
			}
			_conn.Close()
		}(conn)
	}
}

func (s *rpcServer) handler(conner *codec.Conner) error {
	req := &runner.Message{}
	err := conner.Decode(req)
	if err != nil {
		return err
	}
	resp := &runner.Message{ Message: req.Message }
	err = conner.Encode(resp)
	if err != nil {
		return err
	}
	return nil
}
