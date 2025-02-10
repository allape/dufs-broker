package nfs

import (
	"github.com/allape/dufs-broker/ipnet"
	"github.com/allape/gogger"
	"github.com/allape/gohtvfs"
	nfs2 "github.com/willscott/go-nfs"
	nfshelper "github.com/willscott/go-nfs/helpers"
	"net"
)

var l = gogger.New("nfs")

func Start(addr string, dufs *gohtvfs.DufsVFS) error {
	handler := nfshelper.NewNullAuthHandler(NewBillyDufs(dufs))
	cacheHandler := nfshelper.NewCachingHandler(handler, 999)

	addrs, err := ipnet.DescriptAddress(addr)
	if err != nil {
		return err
	}

	for _, addr := range addrs {
		go start(addr, cacheHandler)
	}

	return nil
}

func start(addr string, handler nfs2.Handler) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		l.Error().Fatalf("Failed to listen: %v", err)
	}

	l.Info().Printf("Server running at %s\n", listener.Addr())

	err = nfs2.Serve(listener, handler)
	if err != nil {
		l.Error().Fatalf("Failed to serve NFS: %v", err)
	}
}
