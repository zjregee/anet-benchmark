package main

import (
	"anet-benchmark/runner"
	"anet-benchmark/runner/svr"
)

func NewServer(mode runner.Mode, network, address string) runner.Server {
	switch mode {
	case runner.MODE_RPC:
		return NewRPCServer(network, address)
	case runner.MODE_MUX:
		return NewMuxServer(network, address)
	default:
		panic("illegal value here")
	}
}

func main() {
	svr.Serve(NewServer)
}
