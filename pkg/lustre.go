package pkg

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
	"os"
	"os/exec"
	"strings"
	mount "k8s.io/mount-utils"
)

const (
	nsEnterArg = "--mount=/host/proc/1/ns/mnt"
	// Address of the NFS server
	paramServer = "server"
	// Base directory of the NFS server to create volumes under.
	// The base directory must be a direct child of the root directory.
	// The root directory is omitted from the string, for example:
	//     "base" instead of "/base"
	paramShare            = "share"
	mountOptionsField     = "mountoptions"
	mountPermissionsField = "mountpermissions"
	mountUid              = "uid"
	mountGid              = "gid"
)

type Lustre struct {

}

func IsCorruptedDir(dir string) bool {
	_, pathErr := mount.PathExists(dir)
	return pathErr != nil && mount.IsCorruptedMnt(pathErr)
}

func HostMkdir(pathname string) error {
	var err error
	_, err = exec.Command("nsenter", nsEnterArg, "mkdir", "-p", pathname).Output()
	if err != nil {
		return err
	}
	return nil
}

func HostRmDir(pathname string) error {
	var err error
	_, err = exec.Command("nsenter", nsEnterArg, "rm", "-rf", pathname).Output()
	if err != nil {
		return err
	}
	return nil
}

func HostChmod(pathname string, mod os.FileMode) error {
	var err error
	_, err = exec.Command("nsenter", nsEnterArg, "chmod", fmt.Sprintf("%o", mod), pathname).Output()
	if err != nil {
		return err
	}
	return nil
}

func HostChown(pathname string, uid string, gid string) error {
	var err error
	_, err = exec.Command("nsenter", nsEnterArg, "chown", "-R", fmt.Sprintf("%s:%s", uid, gid), pathname).Output()
	if err != nil {
		return err
	}
	return nil
}

func MountLustre(source string, targetPath string) error {
	var err error
	_, err = exec.Command("nsenter", nsEnterArg, "mount", "-t", "lustre", source, targetPath).Output()
	if err != nil {
		return err
	}
	return nil
}

func EnsureMountLustre(source string, targetPath string) error {
	var count = 5
	cmd := fmt.Sprintf("mkdir -p %s; mount -t lustre %s %s", targetPath, source, targetPath)
	for {
		_, err := exec.Command("nsenter", nsEnterArg, "sh", "-c", cmd).Output()
		if err != nil {
			return nil
		}
		if count = count - 1; count == 0 {
			return status.Errorf(codes.Internal, "EnsureMountLustre error")
		}
	}
}

func HostUMount(targetPath string) error {
	var err error
	_, err = exec.Command("nsenter", nsEnterArg, "umount", targetPath).Output()
	if err != nil {
		return err
	}
	return nil
}

func HostIsFileExist(file string) bool {
	var err error
	_, err = exec.Command("nsenter", nsEnterArg, "ls", file).Output()
	if err != nil {
		klog.V(2).Infof("<<<file not exsit, err is %v>>>", err)
		return false
	}
	return true
}

func IsMountedPoint(targetPath string) bool {
	var mountedPoint []byte
	mountedPoint, _ = exec.Command("nsenter", nsEnterArg, "mount", "-t", "lustre").Output()
	if strings.Index(string(mountedPoint), targetPath) > 0 {
		return true
	}
	return false
}