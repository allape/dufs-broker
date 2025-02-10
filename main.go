package main

import (
	"crypto/tls"
	"github.com/allape/dufs-broker/env"
	"github.com/allape/dufs-broker/ftp"
	"github.com/allape/gogger"
	"github.com/allape/gohtvfs"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

var l = gogger.New("main")

func main() {
	if len(os.Args) >= 2 {
		env.DufsServer = os.Args[1]
	}

	err := gogger.InitFromEnv()
	if err != nil {
		l.Error().Fatalf("Failed to initialize logger: %v", err)
	}

	if env.DufsServer == "" {
		l.Error().Fatalf("Dufs Server Address is required")
	}

	l.Info().Println("Dufs Server:", env.DufsServer)

	u, err := url.Parse(env.DufsServer)
	if err != nil {
		l.Error().Fatalf("Failed to parse DufsServer URL: %v", err)
	}

	caCertPool, err := env.TrustedCertsPoolFromEnv()
	if err != nil {
		l.Error().Fatalf("Failed to create TrustedCertsPool: %v", err)
	}
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	dufs, err := gohtvfs.NewDufsVFS(env.DufsServer)
	if err != nil {
		l.Error().Fatalf("Failed to create DufsVFS: %v", err)
	}
	dufs.HttpClient.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	dufs.SetLogger(gogger.New("dufs").Debug())

	ok, _ := dufs.Online(nil)
	if !ok {
		l.Warn().Println("Dufs server is offline for now")
	} else {
		l.Info().Println("Dufs server is online")
	}

	err = ftp.Start(u, dufs)
	if err != nil {
		l.Error().Fatalf("Failed to start FTP server: %v", err)
	}

	l.Info().Println(env.Banner)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	l.Info().Println("Exiting with", sig)
}
