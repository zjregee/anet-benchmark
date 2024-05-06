package main

import (
	"anet-benchmark/runner"
	"anet-benchmark/runner/svr"
)

func NewServer(network, address string) runner.Server {
	return NewRPCServer(network, address)
}

func main() {
	svr.Serve(NewServer)
}
