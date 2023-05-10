package main

import (
	"flag"
	"toy-lustre-csi/pkg"
)

var (
	endpoint = flag.String("endpoint", "csi.socket", "CSI endpoint")
	nodeId = flag.String("nodeid", "", "CSI nodeid")
)

func main() {
	flag.Parse()
	opt := pkg.DriverOptions{
		Endpoint: *endpoint,
		NodeId: *nodeId,
	}

	d := pkg.NewDriver(opt)
	d.Run()
}

