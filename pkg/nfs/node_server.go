package nfs

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
	"os"
	"path/filepath"
	mount "k8s.io/mount-utils"
	"time"
)

type NodeServer struct {
	nfs *NFSDriver
	mount mount.Interface
}

func (ns *NodeServer) NodeStageVolume(ctx context.Context, request *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return &csi.NodeStageVolumeResponse{}, fmt.Errorf("NodeStageVolume error")
}

func (ns *NodeServer) NodeUnstageVolume(ctx context.Context, request *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return &csi.NodeUnstageVolumeResponse{}, fmt.Errorf("NodeUnstageVolume error")
}

func (ns *NodeServer) NodePublishVolume(ctx context.Context, request *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	klog.Infof("NodePublishVolume Request")
	capacity := request.GetVolumeCapability()
	if capacity == nil {
		return nil, fmt.Errorf("capacity is nil")
	}
	options := capacity.GetMount().GetMountFlags()
	if request.Readonly {
		options = append(options, "ro")
	}

	targetPath := request.GetTargetPath()
	if targetPath == "" {
		return nil, fmt.Errorf("target path nil")
	}

	source := fmt.Sprintf("%s:%s", ns.nfs.nfsServer, filepath.Join(ns.nfs.nfsServerPath, request.GetVolumeId()))

	notMnt, err := ns.mount.IsLikelyNotMountPoint(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(targetPath, os.FileMode(0755)); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			notMnt = true
		}
	} else {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !notMnt {
		return &csi.NodePublishVolumeResponse{}, nil
	}

	klog.Infof("source: %s, targetPath: %s, option: %v", source, targetPath, options)

	if err := ns.mount.Mount(source, targetPath, "nfs", options); err != nil {
		return nil, err
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *NodeServer) NodeUnpublishVolume(ctx context.Context, request *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	klog.Infof("NodeUnpublishVolume Request")
	volumeId := request.GetVolumeId()
	targetPath := request.GetTargetPath()
	klog.Infof("UnpublishVolume targetPath: %s, volumeId: %s", volumeId, targetPath)

	var err error
	extennsiveMountPointCheck := true
	forceUnmounter, ok := ns.mount.(mount.MounterForceUnmounter)
	if ok {
		klog.Infof("force unmout %s on %s", volumeId, targetPath)
		err = mount.CleanupMountWithForce(targetPath, forceUnmounter, extennsiveMountPointCheck, 30*time.Second)
	} else {
		err = mount.CleanupMountPoint(targetPath, ns.mount, true)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmount target %q: %v", targetPath, err)
	}
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (ns *NodeServer) NodeGetVolumeStats(ctx context.Context, request *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return &csi.NodeGetVolumeStatsResponse{}, fmt.Errorf("NodeGetVolumeStats error")
}

func (ns *NodeServer) NodeExpandVolume(ctx context.Context, request *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return &csi.NodeExpandVolumeResponse{}, fmt.Errorf("NodeExpandVolume error")
}

func (ns *NodeServer) NodeGetCapabilities(ctx context.Context, request *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	klog.Infof("NodeGetCapabilities Request")
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: ns.nfs.nodeServiceCapabilities,
	}, nil
}

func (ns *NodeServer) NodeGetInfo(ctx context.Context, request *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	klog.Infof("NodeGetInfo request")
	return &csi.NodeGetInfoResponse{
		NodeId: ns.nfs.nodeId,
	}, nil
}
