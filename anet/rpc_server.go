package main

/*
#cgo LDFLAGS: -luring
#include <liburing.h>
*/
import "C"

import (
	"fmt"
	"net"
	"time"
	"strings"
	"syscall"

	"anet-benchmark/runner"
	"anet-benchmark/net/codec"
)

const QUEUE_DEPTH = 128

func NewRPCServer(network, address string) runner.Server {
	return &rpcServer{
		network: network,
		address: address,
	}
}

var _ runner.Server = &rpcServer{}

type rpcServer struct {
	network string
	address string
}

func (s *rpcServer) Run() error {
	listener, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	defer listener.Close()
	file, err := listener.(*net.TCPListener).File()
	if err != nil {
		return err
	}
	fd := file.Fd()
	var ring C.struct_io_uring
	var sockaddr C.struct_sockaddr
	var addrlen C.socklen_t = C.sizeof_struct_sockaddr
	C.io_uring_queue_init(C.uint(QUEUE_DEPTH), &ring, 0)
	for {
		sqe := C.io_uring_get_sqe(&ring)
		if sqe == nil {
			fmt.Println("an error occurred while accepting, wait 10ms to retry")
			time.Sleep(10 * time.Millisecond)
			continue
		}
		C.io_uring_prep_accept(sqe, C.int(fd), &sockaddr, &addrlen, 0)
		if C.io_uring_submit(&ring) < 0 {
			fmt.Println("an error occurred while accepting, wait 10ms to retry")
			time.Sleep(10 * time.Millisecond)
			continue
		}
		var cqe *C.struct_io_uring_cqe
		for {
			C.io_uring_wait_cqe(&ring, &cqe)
			if cqe.res < 0 {
				errno := syscall.Errno(-cqe.res)
				if errno == syscall.EAGAIN {
					continue
				} else {
					return errno
				}
			} else {
				break
			}
		}
		C.io_uring_cqe_seen(&ring, cqe)
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
