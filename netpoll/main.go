package main

import (
	"anet-benchmark/runner"
	"anet-benchmark/runner/svr"
	"github.com/cloudwego/netpoll"
)

func init() {
	netpoll.DisableGopool()
}

func NewServer(mode runner.Mode, network string, address string) runner.Server {
	switch mode {
	case runner.MODE_RPC:
		return NewRPCServer(network, address)
	default:
		panic("illegal value here")
	}
}

func main() {
	svr.Serve(NewServer)
}
