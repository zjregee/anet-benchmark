package main

import (
	"anet-benchmark/runner"
	"anet-benchmark/runner/cli"
)

func NewClient(mode runner.Mode, network, address string) runner.Client {
	switch mode {
	case runner.MODE_RPC:
		return NewRPCClient(network, address)
	case runner.MODE_MUX:
		return NewMuxClient(network, address, 4)
	default:
		panic("illegal value here")
	}
}

func main() {
	cli.Benching(NewClient)
}
