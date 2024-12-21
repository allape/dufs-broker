package env

import (
	"crypto/x509"
	"github.com/allape/goenv"
	"github.com/allape/gogger"
	"os"
	"strings"
)

const Banner = `
 ______   __   __  _______  _______    _______  ______    _______  ___   _  _______  ______   
|      | |  | |  ||       ||       |  |  _    ||    _ |  |       ||   | | ||       ||    _ |  
|  _    ||  | |  ||    ___||  _____|  | |_|   ||   | ||  |   _   ||   |_| ||    ___||   | ||  
| | |   ||  |_|  ||   |___ | |_____   |       ||   |_||_ |  | |  ||      _||   |___ |   |_||_ 
| |_|   ||       ||    ___||_____  |  |  _   | |    __  ||  |_|  ||     |_ |    ___||    __  |
|       ||       ||   |     _____| |  | |_|   ||   |  | ||       ||    _  ||   |___ |   |  | |
|______| |_______||___|    |_______|  |_______||___|  |_||_______||___| |_||_______||___|  |_|
`

var l = gogger.New("env")

const (
	DubrokerDufsServer   = "DUBROKER_DUFS_SERVER"
	DubrokerTrustedCerts = "DUBROKER_TRUSTED_CERTS"
	DubrokerAddr         = "DUBROKER_ADDRESS"

	DubrokerTlsCertCrt = "DUBROKER_TLS_CERT_CRT"
	DubrokerTlsCertKey = "DUBROKER_TLS_CERT_KEY"

	DubrokerFTPTransferPortRange = "DUBROKER_FTP_TRANSFER_PORT_RANGE"
)

var (
	DufsServer   = goenv.Getenv(DubrokerDufsServer, "http://localhost:5000")
	TrustedCerts = goenv.Getenv(DubrokerTrustedCerts, "")

	//Addr         = goenv.Getenv(DubrokerAddr, "127.0.0.1:2049") // nfs
	//Addr         = goenv.Getenv(DubrokerAddr, "127.0.0.1:2022") // sftp
	Addr = goenv.Getenv(DubrokerAddr, "127.0.0.1:2021")

	TlsCertCrt = goenv.Getenv(DubrokerTlsCertCrt, "") // warn: VLC does not support TLS
	TlsCertKey = goenv.Getenv(DubrokerTlsCertKey, "")

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
