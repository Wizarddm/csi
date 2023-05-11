package pkg

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
	"path/filepath"
	"regexp"
	"strings"
)

type nfsVolume struct {
	// Volume id
	id string
	// Address of the NFS server.
	// Matches paramServer.
	server string
	// Base directory of the NFS server to create volumes under
	// Matches paramShare.
	baseDir string
	// Subdirectory of the NFS server to create volumes under
	subDir string
	// size of volume
	size int64
	// LusterDriver
	ld *LusterDriver
}

func newVolume(name string, size int64, params map[string]string, ld *LusterDriver) (*nfsVolume, error) {
	var server, baseDir string

	for k, v := range params {
		switch strings.ToLower(k) {
		case paramServer:
			server = v
		case paramShare:
			baseDir = v
			if !strings.HasPrefix(baseDir, "/") {
				baseDir = "/" + baseDir
			}
		}
	}
	if server == "" {
		return nil, fmt.Errorf("%v is a required parameter", paramServer)
	}
	if baseDir == "" {
		return nil, fmt.Errorf("%v is a required parameter", paramShare)
	}
	vol := &nfsVolume{
		server: server,
		baseDir: baseDir,
		subDir: name,
		size: size,
		ld: ld,
	}
	vol.id = vol.getVolumeIDFromNfsVol()

	return vol, nil
}

func getNfsVolFromID(id string) (*nfsVolume, error) {
	var server, baseDir, subDir string
	segments := strings.Split(id, separator)
	if len(segments) < 3 {
		klog.Infof("could not split %s into server, baseDir and subDir with separator(%s)", id, separator)
		// try with separator "/"
		volRegex := regexp.MustCompile("^([^/]+)/(.*)/([^/]+)$")
		tokens := volRegex.FindStringSubmatch(id)
		if tokens == nil {
			return nil, fmt.Errorf("could not split %s into server, baseDir and subDir with separator(%s)", id, "/")
		}
		server = tokens[1]
		baseDir = "/" + tokens[2]
		subDir = tokens[3]
	} else {
		server = segments[0]
		baseDir = "/" + segments[1]
		subDir = segments[2]
	}
	return &nfsVolume{
		id:      id,
		server:  server,
		baseDir: baseDir,
		subDir:  subDir,
	}, nil
}

func (vol *nfsVolume) getVolumeIDFromNfsVol() string {
	idElements := make([]string, totalIDElements)
	idElements[idServer] = strings.Trim(vol.server, "/")
	idElements[idBaseDir] = strings.Trim(vol.baseDir, "/")
	idElements[idSubDir] = strings.Trim(vol.subDir, "/")
	return strings.Join(idElements, separator)
}

func (vol *nfsVolume) internalMount() error {
	targetPath := vol.getInternalMountPath()
	if !IsMountedPoint(targetPath) {
		//如果basedir未创建，则先创建
		if !HostIsFileExist(targetPath) {
			klog.Infof("targetPath %v not exist, mkdir", targetPath)
			if err := HostMkdir(targetPath); err != nil {
				klog.Infof("mkdir %v error >>", targetPath)
				return status.Errorf(codes.Internal, fmt.Sprintf("mkdir %v error", targetPath))
			}
		}
		// server:/basedir挂载到workdir
		source := fmt.Sprintf("%s:%s", vol.server, vol.baseDir)
		if err := MountLustre(source, targetPath); err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("mount %v to %v error", source, targetPath))
		}
	}
	return nil
}

func (vol *nfsVolume) internalUnmount() error {
	targetPath := vol.getInternalMountPath()
	if HostIsFileExist(targetPath) {
		if IsMountedPoint(targetPath) {
			if err := HostUMount(targetPath); err != nil {
				return status.Errorf(codes.Internal, fmt.Sprintf("unmont %v error", targetPath))
			}
		}
		HostRmDir(targetPath)
	}
	return nil
}

func (vol *nfsVolume) getInternalMountPath() string {
	return filepath.Join(vol.ld.WorkingMountDir, vol.id)
}

func (vol *nfsVolume) getVolumeSharePath() string {
	return filepath.Join(string(filepath.Separator), vol.baseDir, vol.subDir)
}