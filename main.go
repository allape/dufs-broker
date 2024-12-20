package main

import (
	"crypto/tls"
	vfs "github.com/allape/go-http-vfs"
	"github.com/allape/gogger"
	"net/http"
	"net/url"
	"os"
	"os/signal"
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

	u, err := url.Parse(DufsServer)
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	l.Info().Println("Exiting with", sig)
}
