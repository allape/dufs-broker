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

const Banner = `
 ______   __   __  _______  _______    _______  ______    _______  ___   _  _______  ______   
|      | |  | |  ||       ||       |  |  _    ||    _ |  |       ||   | | ||       ||    _ |  
|  _    ||  | |  ||    ___||  _____|  | |_|   ||   | ||  |   _   ||   |_| ||    ___||   | ||  
| | |   ||  |_|  ||   |___ | |_____   |       ||   |_||_ |  | |  ||      _||   |___ |   |_||_ 
| |_|   ||       ||    ___||_____  |  |  _   | |    __  ||  |_|  ||     |_ |    ___||    __  |
|       ||       ||   |     _____| |  | |_|   ||   |  | ||       ||    _  ||   |___ |   |  | |
|______| |_______||___|    |_______|  |_______||___|  |_||_______||___| |_||_______||___|  |_|
`

const Name = "DUFS FTP Server"

var l = gogger.New("ftp")

func Start(addr string, u *url.URL, dufs *vfs.DufsVFS) error {
	addrs, err := ipnet.DescriptAddress(addr)
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
	return &ftpserver.Settings{
		ListenAddr: d.addr,
		Banner:     Banner,
		PassiveTransferPortRange: &ftpserver.PortRange{
			Start: env.FTPTransferPort - 1,
			End:   env.FTPTransferPort,
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
	return nil, nil
}
