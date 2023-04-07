package nfs

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
)

type Driver struct {
	name     string
	nodeId   string
	version  string
	endpoint string

	// id *identityServer
	ns *csi.NodeServer
}

//func NewDriver() *Driver {
//	return &Driver{}
//}
//
//func (d *Driver) Run(endpoint string, nodeId string)  {
//	server := NewNonBlockGRPCServer()
//
//	nfs := NewNFS("nfs", "1", nodeId)
//
//	ids := &IdentityServer{ n: nfs }
//	cs := &ControllerServer{}
//	ns := &NodeServer{ nfs: nfs }
//
//	server.Start(endpoint, ids, cs, ns)
//}