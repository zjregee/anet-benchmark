package main

import (
	"io"
	"fmt"
	"net"
	"time"
	"strings"

	"anet-benchmark/runner"
	"anet-benchmark/anet/codec"
)

func NewMuxServer(network, address string) runner.Server {
	return &muxServer{
		network: network,
		address: address,
	}
}

var _ runner.Server = &muxServer{}

type muxServer struct {
	network string
	address string
}

func (s *muxServer) Run() error {
	listener, err := net.Listen(s.network, s.address)
	if err != nil {
		panic("can't failed here")
	}
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
		mc := newMuxConn(conn)
		go mc.loopRead()
		go mc.loopWrite()
	}
}

func newMuxConn(conn net.Conn) *muxConn {
	return &muxConn{
		conn: conn,
		conner: codec.NewConner(conn),
		ech: make(chan string),
		wch: make(chan *runner.Message),
	}
}

type muxConn struct {
	conn   net.Conn
	conner *codec.Conner
	ech    chan string
	wch    chan *runner.Message
}

func (mux *muxConn) loopRead() {
	for {
		req := &runner.Message{}
		err := mux.conner.Decode(req)
		if err != nil {
			if err == io.EOF {
				mux.ech <- "EOF"
				return
			}
			panic("can't failed here")
		}
		go func() {
			resp := &runner.Message{ Message: req.Message }
			mux.wch <- resp
		}()
	}
}

func (mux *muxConn) loopWrite() {
	for {
		select {
		case msg := <-mux.wch:
			err := mux.conner.Encode(msg)
			if err != nil {
				panic("can't failed here")
			}
		case <-mux.ech:
			return
		}
	}
}
