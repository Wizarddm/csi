package nfs

import (
	"fmt"
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
	klog.Infof("req name: ", request.GetName())
	mountPath := filepath.Join(c.nfs.nfsMountPath, request.GetName())

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
	klog.Infof("volumeID: ", request.GetVolumeId())
	return nil, os.Remove(filepath.Join(c.nfs.nfsMountPath, request.GetVolumeId()))
}

func (c *ControllerServer) ControllerPublishVolume(ctx context.Context, request *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	return &csi.ControllerPublishVolumeResponse{}, fmt.Errorf("ControllerPublishVolume error")
}

func (c *ControllerServer) ControllerUnpublishVolume(ctx context.Context, request *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return &csi.ControllerUnpublishVolumeResponse{}, fmt.Errorf("ControllerUnpublishVolume error")
}

func (c *ControllerServer) ValidateVolumeCapabilities(ctx context.Context, request *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	return &csi.ValidateVolumeCapabilitiesResponse{}, fmt.Errorf("ValidateVolumeCapabilities error")
}

func (c *ControllerServer) ListVolumes(ctx context.Context, request *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return &csi.ListVolumesResponse{}, fmt.Errorf("ListVolumes error")
}

func (c *ControllerServer) GetCapacity(ctx context.Context, request *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return &csi.GetCapacityResponse{}, fmt.Errorf("GetCapacity error")
}

func (c *ControllerServer) ControllerGetCapabilities(ctx context.Context, request *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	klog.Infof("ControllerGetCapabilities request")
	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: c.nfs.controllerServiceCapabilities,
	}, nil
}

func (c *ControllerServer) CreateSnapshot(ctx context.Context, request *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return &csi.CreateSnapshotResponse{}, fmt.Errorf("CreateSnapshot error")
}

func (c *ControllerServer) DeleteSnapshot(ctx context.Context, request *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return &csi.DeleteSnapshotResponse{}, fmt.Errorf("DeleteSnapshot error")
}

func (c *ControllerServer) ListSnapshots(ctx context.Context, request *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return &csi.ListSnapshotsResponse{}, fmt.Errorf("ListSnapshots error")
}

func (c *ControllerServer) ControllerExpandVolume(ctx context.Context, request *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return &csi.ControllerExpandVolumeResponse{}, fmt.Errorf("ControllerExpandVolume error")
}

func (c *ControllerServer) ControllerGetVolume(ctx context.Context, request *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return &csi.ControllerGetVolumeResponse{}, fmt.Errorf("ControllerGetVolume error")
}


