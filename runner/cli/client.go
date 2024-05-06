package cli

import (
	"flag"

	"anet-benchmark/runner"
)

func Benching(newer runner.ClientNewer) {
	initFlags()
	client := newer(network, address)
	r := runner.NewRunner()
	r.Run(name, client.Echo, concurrent, total, echoSize)
}

var (
	name       string
	address    string
	concurrent int
	total      int64
	echoSize   int

	network    string         = "tcp"
)

func initFlags() {
	flag.StringVar(&name, "name", "", "")
	flag.StringVar(&address, "addr", ":8000", "")
	flag.IntVar(&concurrent, "c", 1, "")
	flag.Int64Var(&total, "n", 1000, "")
	flag.IntVar(&echoSize, "b", 1024, "")
	flag.Parse()
}
