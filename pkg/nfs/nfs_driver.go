package nfs

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/mount-utils"
)

type DriverOptions struct {
	Name            string
	Version         string
	NodeId          string
	Endpoint        string
	Server          string
	ServerPath      string
	WorkingMountDir string
}

type NFSDriver struct {
	name                          string
	version                       string
	nodeId                        string
	endpoint                      string
	nfsMountPath                  string
	nfsServer                     string
	nfsServerPath                 string
	controllerServiceCapabilities []*csi.ControllerServiceCapability
	nodeServiceCapabilities       []*csi.NodeServiceCapability
}

func NewNFSDriver(options *DriverOptions) *NFSDriver {
	nfs := &NFSDriver{
		name:          options.Name,
		version:       options.Version,
		nodeId:        options.NodeId,
		endpoint:      options.Endpoint,
		nfsServer:     options.Server,
		nfsServerPath: options.ServerPath,
		nfsMountPath:  options.WorkingMountDir,
	}

	nfs.addControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
	})

	nfs.addNodeServiceCapabilities([]csi.NodeServiceCapability_RPC_Type{
		//csi.NodeServiceCapability_RPC_SINGLE_NODE_MULTI_WRITER,
		//csi.NodeServiceCapability_RPC_UNKNOWN,
	})

	return nfs
}

func (nfs *NFSDriver) Run()  {
	server := NewNonBlockGRPCServer()
	mounter := mount.New("")
	ids := &IdentityServer{ nfs: nfs }
	cs := &ControllerServer{ nfs: nfs }
	ns := &NodeServer{ nfs: nfs, mount: mounter}

	server.Start(nfs.endpoint, ids, cs, ns)
}

func newControllerServiceCapability(cap csi.ControllerServiceCapability_RPC_Type) *csi.ControllerServiceCapability {
	return &csi.ControllerServiceCapability{
		Type: &csi.ControllerServiceCapability_Rpc{
			Rpc: &csi.ControllerServiceCapability_RPC{
				Type: cap,
			},
		},
	}
}

func (nfs *NFSDriver) addControllerServiceCapabilities(capabilities []csi.ControllerServiceCapability_RPC_Type)  {
	var csc = make([]*csi.ControllerServiceCapability, 0, len(capabilities))
	for _, c := range capabilities {
		csc = append(csc, newControllerServiceCapability(c))
	}
	nfs.controllerServiceCapabilities = csc
}

func newNodeServiceCapability(cap csi.NodeServiceCapability_RPC_Type) *csi.NodeServiceCapability {
	return &csi.NodeServiceCapability{
		Type: &csi.NodeServiceCapability_Rpc{
			Rpc: &csi.NodeServiceCapability_RPC{
				Type: cap,
			},
		},
	}
}

func (nfs *NFSDriver) addNodeServiceCapabilities(capabilities []csi.NodeServiceCapability_RPC_Type) {
	var nsc = make([]*csi.NodeServiceCapability, 0, len(capabilities))
	for _, n := range capabilities {
		nsc = append(nsc, newNodeServiceCapability(n))
	}
	nfs.nodeServiceCapabilities = nsc
}