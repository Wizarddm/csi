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