package nfs

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/net/context"
)

type IdentityServer struct {
	nfs *NFSDriver
}

func (id *IdentityServer) GetPluginInfo(ctx context.Context, request *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	return &csi.GetPluginInfoResponse{
		Name: id.nfs.Name,
		VendorVersion: id.nfs.Version,
	}, nil
}

func (id *IdentityServer) GetPluginCapabilities(ctx context.Context, request *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	return &csi.GetPluginCapabilitiesResponse{}, fmt.Errorf("GetPluginCapabilities error")
}

func (id *IdentityServer) Probe(ctx context.Context, request *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	return nil, fmt.Errorf("err")
}
