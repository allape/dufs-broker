package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	vfs "github.com/allape/go-http-vfs"
	"github.com/allape/gogger"
	"github.com/go-git/go-billy/v5"
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

func NewBillyDufs(dufs *vfs.DufsVFS) billy.Filesystem {
	return &BillyDufs{
		dufs: dufs,
	}
}

var NotImplError = errors.New("not implemented")

type BillyDufs struct {
	billy.Filesystem
	dufs *vfs.DufsVFS
}

// region Basic

func (d BillyDufs) Create(filename string) (billy.File, error) {
	return d.Open(filename)
}

func (d BillyDufs) Open(filename string) (billy.File, error) {
	file, err := d.dufs.Open(filename)
	return &BillyDufsFile{file: file.(*vfs.DufsFile)}, err
}

func (d BillyDufs) OpenFile(filename string, _ int, _ os.FileMode) (billy.File, error) {
	return d.Open(filename)
}

func (d BillyDufs) Stat(filename string) (os.FileInfo, error) {
	return d.dufs.Stat(filename)
}

func (d BillyDufs) Rename(oldpath, newpath string) error {
	return d.dufs.Rename(oldpath, newpath)
}

func (d BillyDufs) Remove(filename string) error {
	return d.dufs.Remove(filename)
}

func (d BillyDufs) Join(elem ...string) string {
	u, err := url.Parse("http://localhost:5000")
	if err != nil {
		panic(err)
	}
	path := u.JoinPath(elem...).Path
	if len(elem) > 1 && strings.HasPrefix(elem[0], "/") && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

// endregion

// region TempFile

func (d BillyDufs) TempFile(_, _ string) (billy.File, error) {
	return nil, NotImplError
}

// endregion

// region Dir

func (d BillyDufs) ReadDir(path string) ([]os.FileInfo, error) {
	dirs, err := d.dufs.ReadDir(path)
	if err != nil {
		return nil, err
	}
	fileInfos := make([]os.FileInfo, len(dirs))
	for i, dir := range dirs {
		info, err := dir.Info()
		if err != nil {
			return nil, err
		}
		fileInfos[i] = info
	}
	return fileInfos, nil
}

func (d BillyDufs) MkdirAll(filename string, perm os.FileMode) error {
	return d.dufs.Mkdir(filename, perm)
}

// endregion

// region Symlink

func (d BillyDufs) Lstat(_ string) (os.FileInfo, error) {
	return nil, NotImplError
}

func (d BillyDufs) Symlink(_, _ string) error {
	return NotImplError
}

func (d BillyDufs) Readlink(_ string) (string, error) {
	return "", NotImplError
}

// endregion

// region Chroot

func (d BillyDufs) Chroot(_ string) (billy.Filesystem, error) {
	return nil, NotImplError
}

func (d BillyDufs) Root() string {
	return d.dufs.Root
}

// endregion

type BillyDufsFile struct {
	billy.File
	file *vfs.DufsFile
}

func (f *BillyDufsFile) Name() string {
	return f.file.Name
}

func (f *BillyDufsFile) Lock() error {
	return NotImplError
}

func (f *BillyDufsFile) Unlock() error {
	return NotImplError
}

func (f *BillyDufsFile) Truncate(_ int64) error {
	return NotImplError
}
