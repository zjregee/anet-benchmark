package svr

import (
	"flag"

	"anet-benchmark/runner"
)

func Serve(newer runner.ServerNewer) {
	initFlags()
	svr := newer(runner.Mode(mode), network, address)
	svr.Run()
}

var (
	mode    int
	address string

	network string = "tcp"
)

func initFlags() {
	flag.IntVar(&mode, "mode", 2, "")
	flag.StringVar(&address, "addr", ":8000", "")
	flag.Parse()
}
