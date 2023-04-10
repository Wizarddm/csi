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
	server := flag.String("server", "","NFS server")
	serverPath := flag.String("serverPath", "", "NFS server path")
	workingMountDir := flag.String("working-mount-dir", "/tmp", "working directory for provisioner to mount nfs shares temporarily")

	klog.Infof("Multi CSI Driver endPoints: %s, nodeId: %s", *endpoint, *nodeId)

	klog.InitFlags(nil)
	defer klog.Flush()
	flag.Parse()

	d := nfs.NewNFSDriver("nfscsi", "1.0", *nodeId, *endpoint, *server, *serverPath, *workingMountDir)
	d.Run()
}