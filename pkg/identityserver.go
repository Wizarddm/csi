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

// GetPluginCapabilities 返回 driver 提供的能力，比如是否提供 controller service 、 volume访问的能力
func (ids *IdentityServer) GetPluginCapabilities(ctx context.Context, request *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	klog.Infof("request GetPluginCapabilities")
	 return &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
		},
	}, nil
}

// Probe 获取CSI插件健康状况
func (ids *IdentityServer) Probe(ctx context.Context, request *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	klog.Infof("request Probe")
	return &csi.ProbeResponse{}, nil
}
