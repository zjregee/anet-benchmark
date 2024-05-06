package main

import (
	"anet-benchmark/runner"
	"anet-benchmark/runner/cli"
)

func NewClient(network, address string) runner.Client {
	return NewRPCClient(network, address)
}

func main() {
	cli.Benching(NewClient)
}
