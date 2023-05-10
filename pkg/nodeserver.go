package pkg

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog/v2"
)

type NodeServer struct {
	lustre *LusterDriver
}

func (ns *NodeServer) NodeStageVolume(ctx context.Context, request *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	panic("implement me")
}

func (ns *NodeServer) NodeUnstageVolume(ctx context.Context, request *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	panic("implement me")
}

func (ns *NodeServer) NodePublishVolume(ctx context.Context, request *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	panic("implement me")
}

func (ns *NodeServer) NodeUnpublishVolume(ctx context.Context, request *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	panic("implement me")
}

func (ns *NodeServer) NodeGetVolumeStats(ctx context.Context, request *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	panic("implement me")
}

func (ns *NodeServer) NodeExpandVolume(ctx context.Context, request *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	panic("implement me")
}

func (ns *NodeServer) NodeGetCapabilities(ctx context.Context, request *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	panic("implement me")
}

func (ns *NodeServer) NodeGetInfo(ctx context.Context, request *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	klog.Infof("NodeGetInfo request")

	return &csi.NodeGetInfoResponse{
		NodeId: ns.lustre.NodeId,
	}, nil
}

