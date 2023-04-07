package nfs

import "github.com/container-storage-interface/spec/lib/go/csi"

type NFSDriver  struct {
	Name                          string
	Version                       string
	NodeId                        string
	Endpoint                      string
	nfsMountPath                  string
	controllerServiceCapabilities []*csi.ControllerServiceCapability
}

func NewNFSDriver(name, version, nodeId, endpoint string) *NFSDriver {
	nfs := &NFSDriver {
		Name: name,
		Version: version,
		NodeId: nodeId,
		Endpoint: endpoint,
	}

	nfs.addControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
	})

	return nfs
}

func (nfs *NFSDriver) Run()  {
	server := NewNonBlockGRPCServer()
	ids := &IdentityServer{ nfs: nfs }
	cs := &ControllerServer{ nfs: nfs }
	ns := &NodeServer{ nfs: nfs }

	server.Start(nfs.Endpoint, ids, cs, ns)
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
