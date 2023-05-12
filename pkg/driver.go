package pkg

import (
	"flag"
	"github.com/container-storage-interface/spec/lib/go/csi"
)

var DefaultUid string
var DefaultGid string
var MountPermissions uint64
var WorkingMountDir string

func init() {
	flag.StringVar(&DefaultUid, "lustre-defaultuid", "nfsnobody", "lustre mount default uid")
	flag.StringVar(&DefaultGid, "lustre-defaultgid", "nfsnobody", "lustre mount default gid")
	flag.Uint64Var(&MountPermissions, "lustre-permissions", 0777, "lustre mount permissions")
	flag.StringVar(&WorkingMountDir, "lustre-mountdir", "/mnt/lustre", "lustre working mount dir")
}

type Driver interface {
	Run()
}

type DriverOptions struct {
	Endpoint string
	NodeId string
	Name string
	Version string
}

func NewDriver(opt DriverOptions) Driver {
	d := &LusterDriver{
		endpoint: opt.Endpoint,
		NodeId: opt.NodeId,
		Name: opt.Name,
		Version: opt.Version,
	}
	d.addControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
	})
	return d
}

type LusterDriver struct {
	Name string
	Version string
	NodeId string
	endpoint string

	DefaultUid       string
	DefaultGid       string
	MountPermissions uint64
	WorkingMountDir  string

	controllerServiceCapabilities []*csi.ControllerServiceCapability
	nodeServiceCapabilities []*csi.NodeServiceCapability
}

func (ld *LusterDriver) Run() {
	s := NewServer()
	ids := &IdentityServer{lustre: ld}
	cs := &ControllerServer{lustre: ld}
	ns := &NodeServer{lustre: ld}
	s.run(ld.endpoint, ids, cs, ns)
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

func (ld *LusterDriver) addControllerServiceCapabilities(capabilities []csi.ControllerServiceCapability_RPC_Type)  {
	var csc = make([]*csi.ControllerServiceCapability, 0, len(capabilities))
	for _, c := range capabilities {
		csc = append(csc, newControllerServiceCapability(c))
	}
	ld.controllerServiceCapabilities = csc
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

func (ld *LusterDriver) addNodeServiceCapabilities(capabilities []csi.NodeServiceCapability_RPC_Type) {
	var nsc = make([]*csi.NodeServiceCapability, 0, len(capabilities))
	for _, n := range capabilities {
		nsc = append(nsc, newNodeServiceCapability(n))
	}
	ld.nodeServiceCapabilities = nsc
}