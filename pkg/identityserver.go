package pkg

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog/v2"
)

type IdentityServer struct {
	lustre *LusterDriver
}

func (ids *IdentityServer) GetPluginInfo(ctx context.Context, request *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	klog.Infof("GetPluginInfo request")
	return &csi.GetPluginInfoResponse{
		Name: ids.lustre.Name,
		VendorVersion: ids.lustre.Version,
	}, nil
}

func (ids *IdentityServer) GetPluginCapabilities(ctx context.Context, request *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	panic("implement me")
}

func (ids *IdentityServer) Probe(ctx context.Context, request *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	panic("implement me")
}
