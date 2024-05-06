package main

import (
	"net"
	"time"

	"anet-benchmark/runner"
	"anet-benchmark/net/codec"
)

func NewRPCClient(network, address string) runner.Client {
	return &rpcClient{
		network: network,
		address: address,
	}
}

var _ runner.Client = &rpcClient{}

type rpcClient struct {
	network string
	address string
}

func (c *rpcClient) Echo(req *runner.Message) (*runner.Message, error) {
	conn, err := net.DialTimeout(c.network, c.address, time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conner := codec.NewConner(conn)
	defer codec.PutConner(conner)
	err = conner.Encode(req)
	if err != nil {
		return nil, err
	}
	resp := &runner.Message{}
	err = conner.Decode(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
