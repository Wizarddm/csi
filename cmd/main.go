package main

import (
	"flag"
	"k8s.io/klog/v2"
	"toy-csi/pkg/nfs"
)

func main() {
	klog.Infof("nfs-csi start...")
	endpoint := flag.String("endpoint", "csi.sock", "CSI endpoint")
	nodeId := flag.String("nodeid", "", "CSI nodeid")
	klog.Infof("Multi CSI Driver endPoints: %s, nodeId: %s", *endpoint, *nodeId)

	klog.InitFlags(nil)
	defer klog.Flush()
	flag.Parse()

	d := nfs.NewNFSDriver("nfs-csi", "1.0", *nodeId, *endpoint)
	d.Run(*endpoint, *nodeId)
}