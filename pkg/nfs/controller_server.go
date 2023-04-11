package nfs

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/net/context"
	"k8s.io/klog/v2"
	"os"
	"path/filepath"
)

type ControllerServer struct {
	nfs *NFSDriver
}

func (c *ControllerServer) CreateVolume(ctx context.Context, request *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	klog.Infof("CreateVolume request")
	klog.Infof("req name: %s", request.GetName())
	mountPath := filepath.Join(c.nfs.nfsMountPath, request.GetName())

	klog.Infof("mkdir %s", mountPath)
	if err := os.Mkdir(mountPath, 755); err != nil {
		klog.Fatalf("mkdir %s error: %s", mountPath, err.Error())
		return nil, err
	}
	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      request.Name,
			CapacityBytes: 0,
		},
	}, nil
}

func (c *ControllerServer) DeleteVolume(ctx context.Context, request *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	klog.Infof("DeleteVolume request")
	klog.Infof("DeleteVolume volumeID: %s", request.GetVolumeId())

	if err := os.Remove(filepath.Join(c.nfs.nfsMountPath, request.GetVolumeId())); err != nil {
		klog.Fatalf("DeleteVolume Failed, DeleteVolume vo: %DeleteVolume: %s", request.GetVolumeId())
		return nil, os.Remove(filepath.Join(c.nfs.nfsMountPath, request.GetVolumeId()))
	}
	return &csi.DeleteVolumeResponse{}, nil
}

func (c *ControllerServer) ControllerPublishVolume(ctx context.Context, request *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	klog.Infof("ControllerPublishVolume request")
	return &csi.ControllerPublishVolumeResponse{}, nil
}

func (c *ControllerServer) ControllerUnpublishVolume(ctx context.Context, request *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	klog.Infof("ControllerUnpublishVolume request")
	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

func (c *ControllerServer) ValidateVolumeCapabilities(ctx context.Context, request *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	klog.Infof("ValidateVolumeCapabilities request")
	return &csi.ValidateVolumeCapabilitiesResponse{}, nil
}

func (c *ControllerServer) ListVolumes(ctx context.Context, request *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	klog.Infof("ListVolumes request")
	return &csi.ListVolumesResponse{}, nil
}

func (c *ControllerServer) GetCapacity(ctx context.Context, request *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	klog.Infof("GetCapacity request")
	return &csi.GetCapacityResponse{}, nil
}

func (c *ControllerServer) ControllerGetCapabilities(ctx context.Context, request *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	klog.Infof("ControllerGetCapabilities request")
	klog.Infof("Capabilities %v", c.nfs.controllerServiceCapabilities)
	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: c.nfs.controllerServiceCapabilities,
	}, nil
}

func (c *ControllerServer) CreateSnapshot(ctx context.Context, request *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	klog.Infof("CreateSnapshot request")
	return &csi.CreateSnapshotResponse{}, nil
}

func (c *ControllerServer) DeleteSnapshot(ctx context.Context, request *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	klog.Infof("DeleteSnapshot request")
	return &csi.DeleteSnapshotResponse{}, nil
}

func (c *ControllerServer) ListSnapshots(ctx context.Context, request *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	klog.Infof("ListSnapshots request")
	return &csi.ListSnapshotsResponse{}, nil
}

func (c *ControllerServer) ControllerExpandVolume(ctx context.Context, request *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	klog.Infof("ControllerExpandVolume request")
	return &csi.ControllerExpandVolumeResponse{}, nil
}

func (c *ControllerServer) ControllerGetVolume(ctx context.Context, request *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	klog.Infof("ControllerGetVolume request")
	return &csi.ControllerGetVolumeResponse{}, nil
}


