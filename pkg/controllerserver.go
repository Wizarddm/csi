package pkg

import (
	"context"
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
	"strconv"
)

type ControllerServer struct {
	lustre *LusterDriver
}

const (
	idServer = iota
	idBaseDir
	idSubDir
	totalIDElements // Always last
)

const separator = "#"

func (cs *ControllerServer) CreateVolume(ctx context.Context, request *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	klog.Infof("request CreateVolume")
	name := request.GetName()
	if len(name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "CreateVolume name must be provided")
	}
	//if err := cs.validateVolumeCapabilities(request.GetVolumeCapabilities()); err != nil {
	//
	//}
	mountPermissions := cs.lustre.MountPermissions
	reqCapacity := request.GetCapacityRange().GetRequiredBytes()
	parameters := request.GetParameters()
	uid := cs.lustre.DefaultUid
	gid := cs.lustre.DefaultGid
	if parameters == nil {
		parameters = make(map[string]string)
	}
	for k, v := range parameters {
		switch k {
		case paramServer:
		// no op
		case paramShare:
		// no op
		case mountOptionsField:
			if v != "" {
				var err error
				if mountPermissions, err = strconv.ParseUint(v, 8, 32); err != nil {
					return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid mountPermissions %s in storage class", v))
				}
			}
		case mountUid:
			uid = v
		case mountGid:
			gid = v
		default:
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid parameter %q in storage class", k))
		}
	}
	vol, err := newVolume(name, reqCapacity, parameters, cs.lustre)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// 把根目录挂载以便创建volume
	if err = vol.internalMount(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mount nfs server: %v", err.Error())
	}
	
	defer func() {
		if err = vol.internalUnmount(); err != nil {
			klog.Warningf("failed to unmount nfs server: %v", err.Error())
		}
	}()

	// create subdirectory under base-dir
	internalVolumePath := vol.getInternalMountPath()
	if HostIsFileExist(internalVolumePath) {
		klog.Warningf("volume %v already exit", vol.id)
	} else if err = HostMkdir(internalVolumePath); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to make subdirectory: %v", err.Error())
	}

	klog.Infof("chown %s with permissions(0%o), uid: %s, gid: %s", name, mountPermissions, uid, gid)
	if err := HostChown(internalVolumePath, uid, gid); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	parameters[paramServer] = vol.server
	parameters[paramShare] = vol.getVolumeSharePath()

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId: vol.id,
			CapacityBytes: 0,
			VolumeContext: parameters,
		},
	}, nil
}

func (cs *ControllerServer) DeleteVolume(ctx context.Context, request *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	volumeId := request.GetVolumeId()
	if volumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "volume id is empty")
	}
	vol, err := getNfsVolFromID(volumeId)
	if err != nil {
		klog.Warningf("failed to get nfs volume for volume id %v deletion: %v", volumeId, err)
		return &csi.DeleteVolumeResponse{}, nil
	}

	// mount nfs base share so we can delete the subdirectory
	if err = vol.internalMount(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mount nfs server: %v", err.Error())
	}
	defer func() {
		if err = vol.internalUnmount(); err != nil {
			klog.Warningf("failed to unmount nfs server: %v", err.Error())
		}
	}()
	// delete subdirectory under base-dir
	internalVolumePath := vol.getInternalMountPath()

	klog.Infof("Removing subdirectory at %v", internalVolumePath)
	if err = HostRmDir(internalVolumePath); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete subdirectory: %v", err.Error())
	}

	return &csi.DeleteVolumeResponse{}, nil
}

func (cs *ControllerServer) ControllerPublishVolume(ctx context.Context, request *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	panic("implement me")
}

func (cs *ControllerServer) ControllerUnpublishVolume(ctx context.Context, request *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	panic("implement me")
}

func (cs *ControllerServer) ValidateVolumeCapabilities(ctx context.Context, request *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	panic("implement me")
}

func (cs *ControllerServer) ListVolumes(ctx context.Context, request *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	panic("implement me")
}

func (cs *ControllerServer) GetCapacity(ctx context.Context, request *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	panic("implement me")
}

// ControllerGetCapabilities 让调用方知道具备哪些能力
func (cs *ControllerServer) ControllerGetCapabilities(ctx context.Context, request *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	klog.Infof("request ControllerGetCapabilities")
	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: cs.lustre.controllerServiceCapabilities,
	}, nil
}

func (cs *ControllerServer) CreateSnapshot(ctx context.Context, request *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	panic("implement me")
}

func (cs *ControllerServer) DeleteSnapshot(ctx context.Context, request *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	panic("implement me")
}

func (cs *ControllerServer) ListSnapshots(ctx context.Context, request *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	panic("implement me")
}

func (cs *ControllerServer) ControllerExpandVolume(ctx context.Context, request *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	panic("implement me")
}

func (cs *ControllerServer) ControllerGetVolume(ctx context.Context, request *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	panic("implement me")
}
