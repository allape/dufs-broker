package sftp

// https://pkg.go.dev/golang.org/x/crypto/ssh#example-NewServerConn

import (
	_ "embed"
	"fmt"
	"github.com/allape/dufs-broker/ipnet"
	"github.com/allape/gogger"
	"github.com/allape/gohtvfs"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"net/url"
)

var (
	//go:embed id_rsa
	HostKey []byte
)

// ssh-keygen -t rsa -f id_rsa

var l = gogger.New("sftp")

func Start(addr string, u *url.URL, dufs *gohtvfs.DufsVFS) error {
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == u.User.Username() {
				if password, ok := u.User.Password(); ok && string(pass) == password {
					return nil, nil
				}
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}

	private, err := ssh.ParsePrivateKey(HostKey)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	config.AddHostKey(private)

	addrs, err := ipnet.DescriptAddress(addr)
	if err != nil {
		return err
	}

	for _, addr := range addrs {
		go start(addr, config, dufs)
	}

	return nil
}

func start(addr string, config *ssh.ServerConfig, dufs *gohtvfs.DufsVFS) {
	l.Info().Println("Starting SFTP server on", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		l.Error().Fatalf("Failed to listen on %s: %v", addr, err)
	}
	defer func() {
		_ = listener.Close()
	}()

	for {
		nConn, err := listener.Accept()
		if err != nil {
			l.Error().Println("Failed to accept incoming connection:", err)
			continue
		}

		go serve(nConn, config)
	}
}

func serve(nConn net.Conn, config *ssh.ServerConfig) {
	_, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		l.Error().Println("Failed to handshake:", err)
		return
	}

	l.Debug().Println("Handshake successful")

	go ssh.DiscardRequests(reqs)

	for c := range chans {
		l.Debug().Println("Incoming channel: %s", c.ChannelType())

		if c.ChannelType() != "session" {
			_ = c.Reject(ssh.UnknownChannelType, "unknown channel type")
			l.Debug().Println("Unknown channel type: %s", c.ChannelType())
			continue
		}

		channel, requests, err := c.Accept()
		if err != nil {
			l.Error().Println("could not accept channel.", err)
			break
		}

		l.Debug().Println("Channel accepted")

		go func(in <-chan *ssh.Request) {
			for req := range in {
				l.Debug().Println("Request: %v", req.Type)
				ok := false
				switch req.Type {
				case "subsystem":
					l.Debug().Println("Subsystem: %s", req.Payload[4:])
					if string(req.Payload[4:]) == "sftp" {
						ok = true
					}
				}
				l.Debug().Println(" - accepted: %v\n", ok)
				_ = req.Reply(ok, nil)
			}
		}(requests)

		serverOptions := []sftp.ServerOption{
			sftp.WithDebug(l.Debug().Writer()),
		}

		//if readOnly {
		//	serverOptions = append(serverOptions, sftp.ReadOnly())
		//	l.Debug().Println("Read-only server")
		//} else {
		//	l.Debug().Println("Read write server")
		//}

		server, err := sftp.NewServer(
			channel,
			serverOptions...,
		)
		if err != nil {
			l.Error().Println("Failed to create server:", err)
			break
		}
		if err := server.Serve(); err != nil {
			if err != io.EOF {
				l.Error().Println("sftp server completed with error:", err)
				break
			}
		}

		err = server.Close()
		if err != nil {
			l.Error().Println("Failed to close server:", err)
		}

		l.Debug().Println("sftp client exited session.")
	}
}
