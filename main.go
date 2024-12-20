package main

import (
	"crypto/tls"
	"fmt"
	vfs "github.com/allape/go-http-vfs"
	"github.com/allape/gogger"
	"github.com/willscott/go-nfs"
	nfshelper "github.com/willscott/go-nfs/helpers"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var l = gogger.New("main")

func main() {
	if len(os.Args) >= 2 {
		DufsServer = os.Args[1]
	}

	err := gogger.InitFromEnv()
	if err != nil {
		l.Error().Fatalf("Failed to initialize logger: %v", err)
	}

	if DufsServer == "" {
		l.Error().Fatalf("Dufs Server Address is required")
	}

	// TODO auth over NFS
	_, err = url.Parse(DufsServer)
	if err != nil {
		l.Error().Fatalf("Failed to parse DufsServer URL: %v", err)
	}

	caCertPool, err := TrustedCertsPoolFromEnv()
	if err != nil {
		l.Error().Fatalf("Failed to create TrustedCertsPool: %v", err)
	}
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	dufs, err := vfs.NewDufsVFS(DufsServer)
	if err != nil {
		l.Error().Fatalf("Failed to create DufsVFS: %v", err)
	}
	dufs.HttpClient.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	dufs.SetLogger(gogger.New("dufs").Debug())

	handler := nfshelper.NewNullAuthHandler(NewBillyDufs(dufs))
	cacheHandler := nfshelper.NewCachingHandler(handler, 999)

	if strings.HasPrefix(Addr, ":") {
		interfaces, err := net.Interfaces()
		if err != nil {
			l.Error().Fatalf("Failed to get interfaces: %v", err)
		}

		for _, iface := range interfaces {
			addrs, err := iface.Addrs()
			if err != nil {
				l.Warn().Println("Error getting addresses for interface:", iface.Name, err)
				continue
			}

			for _, addr := range addrs {
				ipAddr, ok := addr.(*net.IPNet)
				if !ok {
					continue
				} else if ipAddr.IP.IsMulticast() || ipAddr.IP.IsLinkLocalMulticast() || ipAddr.IP.IsLinkLocalUnicast() {
					continue
				}

				if ipAddr.IP.To16() == nil {
					go start(fmt.Sprintf("%s%s", ipAddr.IP, Addr), cacheHandler)
				} else {
					go start(fmt.Sprintf("[%s]%s", ipAddr.IP, Addr), cacheHandler)
				}
			}
		}
	} else {
		go start(Addr, cacheHandler)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	l.Info().Println("Exiting with", sig)
}

func start(addr string, handler nfs.Handler) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		l.Error().Fatalf("Failed to listen: %v", err)
	}

	l.Info().Printf("Server running at %s\n", listener.Addr())

	err = nfs.Serve(listener, handler)
	if err != nil {
		l.Error().Fatalf("Failed to serve NFS: %v", err)
	}
}
