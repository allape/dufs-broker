package env

import (
	"crypto/x509"
	"github.com/allape/goenv"
	"github.com/allape/gogger"
	"os"
	"strings"
)

var l = gogger.New("env")

const (
	DubrokerDufsServer   = "DUBROKER_DUFS_SERVER"
	DubrokerTrustedCerts = "DUBROKER_TRUSTED_CERTS"
	DubrokerAddr         = "DUBROKER_ADDRESS"

	DubrokerFTPTransferPortRange = "DUBROKER_FTP_TRANSFER_PORT_RANGE"
)

var (
	DufsServer   = goenv.Getenv(DubrokerDufsServer, "http://localhost:5000")
	TrustedCerts = goenv.Getenv(DubrokerTrustedCerts, "")

	//Addr         = goenv.Getenv(DubrokerAddr, "127.0.0.1:2049") // nfs
	//Addr         = goenv.Getenv(DubrokerAddr, "127.0.0.1:2022") // sftp
	Addr = goenv.Getenv(DubrokerAddr, "127.0.0.1:2021")

	FTPTransferPortRange = goenv.Getenv(DubrokerFTPTransferPortRange, PortRange("50000-50100"))
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
