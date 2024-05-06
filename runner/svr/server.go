package svr

import (
	"flag"

	"anet-benchmark/runner"
)

func Serve(newer runner.ServerNewer) {
	initFlags()
	svr := newer(network, address)
	svr.Run()
}

var (
	address string

	network string = "tcp"
)

func initFlags() {
	flag.StringVar(&address, "addr", ":8000", "")
	flag.Parse()
}
