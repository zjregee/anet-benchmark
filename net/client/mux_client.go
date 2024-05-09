package main

import (
	"net"
	"time"
	"sync/atomic"

	"anet-benchmark/runner"
	"anet-benchmark/net/codec"
)

func NewMuxClient(network, address string, size int) runner.Client {
	c := &muxClient{
		network: network,
		address: address,
		size: uint32(size),
		conns: make([]*muxConn, size),
	}
	for i := range c.conns {
		conn, err := net.DialTimeout(network, address, time.Second)
		if err != nil {
			panic("can't failed here")
		}
		mc := newMuxConn(conn)
		go mc.loopRead()
		go mc.loopWrite()
		c.conns[i] = mc
	}
	return c
}

type muxClient struct {
	network string
	address string
	conns   []*muxConn
	size    uint32
	cursor  uint32
}

func (c *muxClient) Echo(req *runner.Message) (*runner.Message, error) {
	conn := c.conns[atomic.AddUint32(&c.cursor, 1) % c.size]
	conn.wch <- req
	resp := <-conn.rch
	return resp, nil
}

func newMuxConn(conn net.Conn) *muxConn {
	return &muxConn{
		conn: conn,
		conner: codec.NewConner(conn),
		rch: make(chan *runner.Message),
		wch: make(chan *runner.Message),
	}
}

type muxConn struct {
	conn   net.Conn
	conner *codec.Conner
	rch chan *runner.Message
	wch chan *runner.Message
}

func (mc *muxConn) loopRead() {
	for {
		msg := &runner.Message{}
		err := mc.conner.Decode(msg)
		if err != nil {
			panic("can't failed here")
		}
		mc.rch <- msg
	}
}

func (mc *muxConn) loopWrite() {
	for {
		msg := <-mc.wch
		err := mc.conner.Encode(msg)
		if err != nil {
			panic("can't failed here")
		}
	}
}
