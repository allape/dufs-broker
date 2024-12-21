package ftp

import (
	"crypto/tls"
	_ "embed"
	"fmt"
	"github.com/allape/dufs-broker/env"
	"github.com/allape/dufs-broker/ipnet"
	vfs "github.com/allape/go-http-vfs"
	"github.com/allape/gogger"
	ftpserver "github.com/fclairamb/ftpserverlib"
	"net/url"
)

const Name = "DUFS FTP Server"

var l = gogger.New("ftp")

func Start(u *url.URL, dufs *vfs.DufsVFS) error {
	addrs, err := ipnet.DescriptAddress(env.Addr)
	if err != nil {
		return err
	}

	for _, addr := range addrs {
		go start(addr, u, dufs)
	}

	return nil
}

func start(addr string, u *url.URL, dufs *vfs.DufsVFS) {
	server := ftpserver.NewFtpServer(&DufsDriver{
		addr: addr,
		u:    u,
		dufs: dufs,
	})
	server.Logger = NewLogger(l)

	l.Info().Println("FTP server started", "addr", addr)

	err := server.ListenAndServe()
	if err != nil {
		l.Error().Println("FTP server error", "error", err)
	}
}

type DufsDriver struct {
	ftpserver.MainDriver
	addr string
	u    *url.URL
	dufs *vfs.DufsVFS
}

func (d *DufsDriver) GetSettings() (*ftpserver.Settings, error) {
	pStart, pEnd, err := env.FTPTransferPortRange.Range()
	if err != nil {
		return nil, err
	}

	tlsMode := ftpserver.MandatoryEncryption

	if env.TlsCertCrt == "" || env.TlsCertKey == "" {
		tlsMode = ftpserver.ClearOrEncrypted
	}

	return &ftpserver.Settings{
		ListenAddr:  d.addr,
		Banner:      env.Banner,
		TLSRequired: tlsMode,
		PassiveTransferPortRange: &ftpserver.PortRange{
			Start: pStart,
			End:   pEnd,
		},
	}, nil
}

func (d *DufsDriver) ClientConnected(cc ftpserver.ClientContext) (string, error) {
	l.Debug().Println("Client connected", cc.ID(), cc.Path())
	return Name, nil
}

func (d *DufsDriver) ClientDisconnected(cc ftpserver.ClientContext) {
	l.Debug().Println("Client disconnected", cc.ID(), cc.Path())
}

func (d *DufsDriver) AuthUser(_ ftpserver.ClientContext, user, pass string) (ftpserver.ClientDriver, error) {
	if d.u.User.Username() == "" {
		return &DufsClientDriver{
			dufs: d.dufs,
		}, nil
	}

	if user == d.u.User.Username() {
		if password, ok := d.u.User.Password(); ok && pass == password {
			return &DufsClientDriver{
				dufs: d.dufs,
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid user or password")
}

func (d *DufsDriver) GetTLSConfig() (*tls.Config, error) {
	var tlsConfig *tls.Config

	if env.TlsCertCrt == "" || env.TlsCertKey == "" {
		return tlsConfig, nil
	}

	l.Info().Println("Using TLS", "cert", env.TlsCertCrt, "key", env.TlsCertKey)

	cert, err := tls.LoadX509KeyPair(env.TlsCertCrt, env.TlsCertKey)
	if err != nil {
		return nil, err
	}

	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	return tlsConfig, nil
}
