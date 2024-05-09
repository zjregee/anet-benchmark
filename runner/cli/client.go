package cli

import (
	"flag"

	"anet-benchmark/runner"
)

func Benching(newer runner.ClientNewer) {
	initFlags()
	client := newer(runner.Mode(mode), network, address)
	r := runner.NewRunner()
	r.Run(name, client.Echo, total, echosize, concurrent)
}

var (
	mode       int
	name       string
	address    string
	total      int
	echosize   int
	concurrent int

	network    string         = "tcp"
)

func initFlags() {
	flag.IntVar(&mode, "mode", 2, "")
	flag.StringVar(&name, "name", "", "")
	flag.StringVar(&address, "addr", ":8000", "")
	flag.IntVar(&total, "total", 10000, "")
	flag.IntVar(&echosize, "echosize", 1024, "")
	flag.IntVar(&concurrent, "concurrent", 4, "")
	flag.Parse()
}
