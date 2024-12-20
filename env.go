package main

import (
	"crypto/x509"
	"github.com/allape/goenv"
	"os"
	"strings"
)

const (
	EnvDufsServer   = "DUBROKER_DUFS_SERVER"
	EnvTrustedCerts = "DUBROKER_TRUSTED_CERTS"
	EnvAddr         = "DUBROKER_ADDRESS"
)

var (
	DufsServer   = goenv.Getenv(EnvDufsServer, "http://localhost:5000")
	TrustedCerts = goenv.Getenv(EnvTrustedCerts, "")
	Addr         = goenv.Getenv(EnvAddr, "127.0.0.1:2022")
)

func TrustedCertsPoolFromEnv() (*x509.CertPool, error) {
	l.Info().Println("TrustedCerts:", TrustedCerts)

	certs := strings.Split(TrustedCerts, ",")

	caCertPool := x509.NewCertPool()

	for _, cert := range certs {
		cert = strings.TrimSpace(cert)
		if cert == "" {
			continue
		}
		caCert, err := os.ReadFile(cert)
		if err != nil {
			return nil, err
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}

	return caCertPool, nil
}
