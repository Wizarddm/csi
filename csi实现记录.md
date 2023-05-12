# 记录一下csi整个实现过程

## 项目初始化

创建项目目录

- cmd

- deploy

- pkg

  - controllerserver.go
  - identityserver.go
  - nodeserver.go

  - driver.go
  - lustre.go
  - utils.go
  - server.go

- build.sh

- Dockerfile

## 实现三个grpc service

实现一个CSI需要实现三个grpc service，这三个grpc service的定义如下：

- **IdentotyServer**

```go
// github.com/container-storage-interface/spec/lib/go/csi/csi.pb.go
type IdentityServer interface {
    //返回driver的信息，比如名字，版本
    GetPluginInfo(context.Context, *GetPluginInfoRequest) (*GetPluginInfoResponse, error)
    //返回driver提供的能力，比如是否提供Controller Service,volume 访问能能力
    GetPluginCapabilities(context.Context, *GetPluginCapabilitiesRequest) (*GetPluginCapabilitiesResponse, error)
    //探针
    Probe(context.Context, *ProbeRequest) (*ProbeResponse, error)
}
```

- **ControllerServer**

```go
// github.com/container-storage-interface/spec/lib/go/csi/csi.pb.go
type ControllerServer interface {
    //创建卷
    CreateVolume(context.Context, *CreateVolumeRequest) (*CreateVolumeResponse, error)
    //删除卷
    DeleteVolume(context.Context, *DeleteVolumeRequest) (*DeleteVolumeResponse, error)
    //attach 卷
    ControllerPublishVolume(context.Context, *ControllerPublishVolumeRequest) (*ControllerPublishVolumeResponse, error)
    //unattach卷
    ControllerUnpublishVolume(context.Context, *ControllerUnpublishVolumeRequest) (*ControllerUnpublishVolumeResponse, error)
    //返回存储卷的功能点，如是否支持挂载到多个节点上，是否支持多个节点同时读写
    ValidateVolumeCapabilities(context.Context, *ValidateVolumeCapabilitiesRequest) (*ValidateVolumeCapabilitiesResponse, error)
    //列出所有卷
    ListVolumes(context.Context, *ListVolumesRequest) (*ListVolumesResponse, error)
    //返回存储资源池的可用空间大小
    GetCapacity(context.Context, *GetCapacityRequest) (*GetCapacityResponse, error)
    //返回controller插件的功能点，如是否支持GetCapacity接口，是否支持snapshot功能等
    ControllerGetCapabilities(context.Context, *ControllerGetCapabilitiesRequest) (*ControllerGetCapabilitiesResponse, error)
    //创建快照
    CreateSnapshot(context.Context, *CreateSnapshotRequest) (*CreateSnapshotResponse, error)
    //删除快照
    DeleteSnapshot(context.Context, *DeleteSnapshotRequest) (*DeleteSnapshotResponse, error)
    //列出快照
    ListSnapshots(context.Context, *ListSnapshotsRequest) (*ListSnapshotsResponse, error)
    //扩容
    ControllerExpandVolume(context.Context, *ControllerExpandVolumeRequest) (*ControllerExpandVolumeResponse, error)
    //获得卷
    ControllerGetVolume(context.Context, *ControllerGetVolumeRequest) (*ControllerGetVolumeResponse, error)
}
```

- **NodeServer**

```go
// github.com/container-storage-interface/spec/lib/go/csi/csi.pb.go
type NodeServer interface {
      //如果存储卷没有格式化，首先要格式化。然后把存储卷mount到一个临时的目录（这个目录通常是节点上的一个全局目录）。再通过NodePublishVolume将存储卷mount到pod的目录中。mount过程分为2步，原因是为了支持多个pod共享同一个volume（如NFS）。
    NodeStageVolume(context.Context, *NodeStageVolumeRequest) (*NodeStageVolumeResponse, error)
    //NodeStageVolume的逆操作，将一个存储卷从临时目录umount掉
    NodeUnstageVolume(context.Context, *NodeUnstageVolumeRequest) (*NodeUnstageVolumeResponse, error)
    //将存储卷从临时目录mount到目标目录（pod目录）
    NodePublishVolume(context.Context, *NodePublishVolumeRequest) (*NodePublishVolumeResponse, error)
    //将存储卷从pod目录umount掉
    NodeUnpublishVolume(context.Context, *NodeUnpublishVolumeRequest) (*NodeUnpublishVolumeResponse, error)
    //返回可用于该卷的卷容量统计信息
    NodeGetVolumeStats(context.Context, *NodeGetVolumeStatsRequest) (*NodeGetVolumeStatsResponse, error)
    //node上执行卷扩容
    NodeExpandVolume(context.Context, *NodeExpandVolumeRequest) (*NodeExpandVolumeResponse, error)
    //返回Node插件的功能点，如是否支持stage/unstage功能
    NodeGetCapabilities(context.Context, *NodeGetCapabilitiesRequest) (*NodeGetCapabilitiesResponse, error)
    //返回节点信息
    NodeGetInfo(context.Context, *NodeGetInfoRequest) (*NodeGetInfoResponse, error)
}
```

这些接口并不是需要全部实现



[容器存储接口](https://github.com/container-storage-interface/spec/blob/master/spec.md)（CSI）是用于将任意块和[文件存储](https://cloud.tencent.com/product/cfs?from=20067&from_column=20067)系统暴露给诸如Kubernetes之类的[容器](https://cloud.tencent.com/product/tke?from=20067&from_column=20067)编排系统（CO）上的容器化工作负载的标准。 使用CSI的第三方存储提供商可以编写和部署在Kubernetes中公开新存储系统的插件，而无需接触核心的Kubernetes代码。

具体来说，Kubernetes针对CSI规定了以下内容：

- Kubelet到CSI驱动程序的通信

  \- Kubelet通过Unix域套接字直接向CSI驱动程序发起CSI调用（例如`NodeStageVolume`，`NodePublishVolume`等），以挂载和卸载卷。

  \- Kubelet通过kubelet插件注册机制发现CSI驱动程序（以及用于与CSI驱动程序进行交互的Unix域套接字）。

  \- 因此，部署在Kubernetes上的所有CSI驱动程序**必须**在每个受支持的节点上使用kubelet插件注册机制进行注册。

- Master到CSI驱动程序的通信

  \- Kubernetes master组件不会直接（通过Unix域套接字或其他方式）与CSI驱动程序通信。

  \- Kubernetes master组件仅与Kubernetes API交互。

  \- 因此，需要依赖于Kubernetes API的操作的CSI驱动程序（例如卷创建，卷attach，卷快照等）必须监听Kubernetes API并针对它触发适当的CSI操作（例如下面的一系列的external组件）。

## 组件

![img](https://static001.geekbang.org/infoq/8b/8b17561f8acee239cbc14a416396280b.png)

CSI实现中的组件分为两部分：

- 由k8s官方维护的一系列external组件负责注册CSI driver 或监听k8s对象资源，从而发起csi driver调用，比如（node-driver-registrar，external-attacher，external-provisioner，external-resizer，external-snapshotter，livenessprobe）
- 各云厂商or开发者自行开发的组件（需要实现CSI Identity，CSI Controller，CSI Node 接口）

### External 组件（k8s Team）

这部分组件是由k8s官方提供的，作为k8s api跟csi driver的桥梁：

- **node-driver-registrar**

  CSI node-driver-registrar是一个sidecar容器，可从CSI driver获取驱动程序信息（使用NodeGetInfo），并使用kubelet插件注册机制在该节点上的kubelet中对其进行注册。

- **external-attacher**

  它是一个sidecar容器，用于监视Kubernetes VolumeAttachment对象并针对驱动程序端点触发CSI ControllerPublish和ControllerUnpublish操作

- **external-provisioner**

  它是一个sidecar容器，用于监视Kubernetes PersistentVolumeClaim对象并针对驱动程序端点触发CSI CreateVolume和DeleteVolume操作。

  external-attacher还支持快照数据源。 如果将快照CRD资源指定为PVC对象上的数据源，则此sidecar容器通过获取SnapshotContent对象获取有关快照的信息，并填充数据源字段，该字段向存储系统指示应使用指定的快照填充新卷 。

- **external-resizer**

  它是一个sidecar容器，用于监视Kubernetes API[服务器](https://cloud.tencent.com/product/cvm?from=20067&from_column=20067)上的PersistentVolumeClaim对象的改动，如果用户请求在PersistentVolumeClaim对象上请求更多存储，则会针对CSI端点触发ControllerExpandVolume操作。

- **external-snapshotter**

  它是一个sidecar容器，用于监视Kubernetes API服务器上的VolumeSnapshot和VolumeSnapshotContent CRD对象。创建新的VolumeSnapshot对象（引用与此驱动程序对应的SnapshotClass CRD对象）将导致sidecar容器提供新的快照。

  该Sidecar侦听指示成功创建VolumeSnapshot的服务，并立即创建VolumeSnapshotContent资源。

- **livenessprobe**

  它是一个sidecar容器，用于监视CSI驱动程序的运行状况，并通过[Liveness Probe机制](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)将其报告给Kubernetes。 这使Kubernetes能够自动检测驱动程序问题并重新启动Pod以尝试解决问题。

- dirver.go

``` go
package pkg

type Driver interface {
	run()
}

type DriverOptions struct {
	endpoint string
	nodeId string
}

func NewDriver(opt DriverOptions) Driver {
	return &lusterDriver{
		endpoint: opt.endpoint,
		nodeId: opt.nodeId,
	}
}

type lusterDriver struct {
	endpoint string
	nodeId string
}

func (ld *lusterDriver) run() {
	s := NewServer()
	// Name和Version暂时固定随意写一个
	lustre := &Lustre{
		NodeId: ld.nodeId,
		Name: "lustre",
		Version: "1.0",
	}
	ids := &IdentityServer{lustre: lustre}
	cs := &ControllerServer{}
	ns := &NodeServer{lustre: lustre}
	s.run(ld.endpoint, ids, cs, ns)
}
```

- server.go

server是实现grpc server

```go
package pkg

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"
	"net"
	"os"
)

type Server interface {
	run(endpoint string, ids csi.IdentityServer, cs csi.ControllerServer, ns csi.NodeServer)
}

func NewServer() Server {
	return &server{}
}

type server struct {

}

func (ld *server) run(endpoint string, ids csi.IdentityServer, cs csi.ControllerServer, ns csi.NodeServer) {
	proto, addr, err := ParseEndpoint(endpoint)
	if err != nil {
		klog.Fatalf("parse endpoint error: %s", err.Error())
	}

	if proto == "unix" {
		addr = "/" + addr
		if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
			klog.Fatalf("Failed to remove %s, error: %s", addr, err.Error())
		}
	}

	listener, err := net.Listen(proto, addr)
	if err != nil {
		klog.Fatalf("Failed to listen, %s", err)
	}
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(logGRPC),
	}

	server := grpc.NewServer(opts...)

	if ids != nil {
		csi.RegisterIdentityServer(server, ids)
	}
	if cs != nil {
		csi.RegisterControllerServer(server, cs)
	}
	if ns != nil {
		csi.RegisterNodeServer(server, ns)
	}
	err = server.Serve(listener)
	if err != nil {
		klog.Fatalf("failed to serve grpc server: %v", err)
	}

}
```

- utils

用于记录公共方法

```go
package pkg

import (
	"context"
	"fmt"
	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"
	"strings"
)

func ParseEndpoint(ep string) (string, string, error) {
	if strings.HasPrefix(strings.ToLower(ep), "unix://") || strings.HasPrefix(strings.ToLower(ep),"tcp://") {
		s := strings.SplitN(ep, "://", 2)
		if s[1] != "" {
			return s[0], s[1], nil
		}
	}
	return "", "", fmt.Errorf("Invalid endpoint: %v", ep)
}

func logGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	level := klog.Level(2)
	klog.V(level).Infof("GRPC call: %s", info.FullMethod)
	klog.V(level).Infof("GRPC request: %s", protosanitizer.StripSecrets(req))

	resp, err := handler(ctx, req)
	if err != nil {
		klog.Errorf("GRPC error: %v", err)
	} else {
		klog.V(level).Infof("GRPC response: %s", protosanitizer.StripSecrets(resp))
	}
	return resp, err
}
```

编译部署

build.sh

```go
#!/bin/bash

set -ex

# 编译
CGO_ENABLED=0 go build -mod=vendor -o build/toy-lustre-csi/ cmd/main.go

image="toy-lustre-csi:1.0"

# 打包镜像
docker build -t $image .

# 推送镜像
# docker push image
```

Dockerfile

```go
FROM busybox

COPY build/toy-lustre-csi /

ENTRYPOINT ["/toy-lustre-csi"]
```



## 实现注册过程

CSI插件注册过程只会调用CSI进程的两个方法，这两个方法分别是`IdentityServer下的GetPluginInfo方法`和`NodeServer下的NodeGetInfo方法`，我们先实现这两个方法验证下验证下注册过程

- **IdentityServer下的GetPluginInfo方法**

```go
package pkg

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog/v2"
)

type IdentityServer struct {
	lustre *Lustre
}

func (ids *IdentityServer) GetPluginInfo(ctx context.Context, request *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	klog.Infof("GetPluginInfo request")
	return &csi.GetPluginInfoResponse{
		Name: ids.lustre.Name,
		VendorVersion: ids.lustre.Version,
	}, nil
}

```

- **NodeServer下的NodeGetInfo方法**

```go
package pkg

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog/v2"
)

type NodeServer struct {
	lustre *Lustre
}

func (ns *NodeServer) NodeGetInfo(ctx context.Context, request *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	klog.Infof("NodeGetInfo request")

	return &csi.NodeGetInfoResponse{
		NodeId: ns.lustre.NodeId,
	}, nil
}
```

##### 注册过程产物

注册过程中会有如下产物：

- node-driver-registrar进程的sock文件：/var/lib/kubelet/plugins_registry/{csiDriverName}-reg.sock
- CSI进程的sock文件：/var/lib/kubelet/plugins/{xxx}/csi.sock
- 节点对应Node对象的annotation中会有一个关于该CSI插件的注解
- 会有一个CSINode对象



用上述脚本把代码编译打包成镜像`toy-lustre-csi:1.0`，之后准备部署的yaml（当前阶段为了简单起见，给daemonSet的nodeSelector加了个`kubernetes.io/hostname`标签用于指定只在一个节点上运行pod）：

```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test-lustre
  namespace: test-z

---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: test-lustre
rules:
  - apiGroups: ["storage.k8s.io"]
    resources: ["csinodes"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "watch"]
---

kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: test-lustre
subjects:
  - kind: ServiceAccount
    name: test-lustre
    namespace: test-z
roleRef:
  kind: ClusterRole
  name: test-lustre
  apiGroup: rbac.authorization.k8s.io

---
kind: DaemonSet
apiVersion: apps/v1
metadata:
  name: test-lustre
  namespace: test-z
spec:
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 1
    type: RollingUpdate
  selector:
    matchLabels:
      app: test-lustre
  template:
    metadata:
      labels:
        app: test-lustre
    spec:
      hostNetwork: true  # original nfs connection would be broken without hostNetwork setting
      dnsPolicy: Default  # available values: Default, ClusterFirstWithHostNet, ClusterFirst
      serviceAccountName: test-lustre
      nodeSelector:
        kubernetes.io/os: linux
        kubernetes.io/hostname: worker3 # 调试阶段可以先指定某一个节点启动
      tolerations:
        - operator: "Exists"
      containers:
        - name: node-driver-registrar
          image: objectscale/csi-node-driver-registrar:v2.5.0
          imagePullPolicy: IfNotPresent
          args:
            - --v=2
            - --csi-address=/csi/csi.sock
            - --kubelet-registration-path=$(DRIVER_REG_SOCK_PATH)
          env:
            - name: DRIVER_REG_SOCK_PATH
              value: /var/lib/kubelet/plugins/csi-nfsplugin/csi.sock
            - name: KUBE_NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
          volumeMounts:
            - name: socket-dir
              mountPath: /csi
            - name: registration-dir
              mountPath: /registration
          resources:
            limits:
              memory: 100Mi
            requests:
              cpu: 10m
              memory: 20Mi
        - name: lustre-csi
          securityContext:
            privileged: true
            capabilities:
              add: ["SYS_ADMIN"]
            allowPrivilegeEscalation: true
          image: 172.31.4.89/test/toy-lustre-csi:1.0
          imagePullPolicy: "Always"
          args:
            - "--endpoint=$(CSI_ENDPOINT)"
            - "--nodeid=$(NODE_ID)"
          env:
            - name: NODE_ID
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: CSI_ENDPOINT
              value: unix://csi/csi.sock
          volumeMounts:
            - name: socket-dir
              mountPath: /csi
          resources:
            limits:
              memory: 300Mi
            requests:
              cpu: 10m
              memory: 20Mi
      volumes:
        - name: socket-dir
          hostPath:
            path: /var/lib/kubelet/plugins/csi-nfsplugin
            type: DirectoryOrCreate
        - name: registration-dir
          hostPath:
            path: /var/lib/kubelet/plugins_registry
            type: Directory
```

## dynamic provisioning

在csi接口中，对应的是`ControllerServer`下的`CreateVolume`和`DeleteVolume`方法，现在来实现这两个方法。在实现这两个方法前需要通过`ControllerServer`的`ControllerGetCapabilities`方法让调用放知道自己有`CreateVolume/DeleteVolume`的能力

```go
func NewDriver(opt DriverOptions) Driver {

	d := &lusterDriver{
		endpoint: opt.Endpoint,
		nodeId: opt.NodeId,
	}
	d.addControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
	})
	return d
}


func newControllerServiceCapability(cap csi.ControllerServiceCapability_RPC_Type) *csi.ControllerServiceCapability {
	return &csi.ControllerServiceCapability{
		Type: &csi.ControllerServiceCapability_Rpc{
			Rpc: &csi.ControllerServiceCapability_RPC{
				Type: cap,
			},
		},
	}
}

func (ld *lusterDriver) addControllerServiceCapabilities(capabilities []csi.ControllerServiceCapability_RPC_Type)  {
	var csc = make([]*csi.ControllerServiceCapability, 0, len(capabilities))
	for _, c := range capabilities {
		csc = append(csc, newControllerServiceCapability(c))
	}
	ld.controllerServiceCapabilities = csc
}

// ControllerGetCapabilities 让调用方知道具备哪些能力
func (cs *ControllerServer) ControllerGetCapabilities(ctx context.Context, request *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	klog.Infof("request ControllerGetCapabilities")
	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: cs.lustre.controllerServiceCapabilities,
	}, nil
}
```

再看看CreateVolume和DeleteVolume方法的实现：

- CreateVolume

``` go
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
```

- DeleteVolume

```go
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
```

额外需要实现`IdentityServer`的`Probe`方法（用来获取csi健康状况）及`GetPluginCapabilities`方法（返回 driver 提供的能力，比如是否提供 controller service 、 volume访问的能力）

```go
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
```

