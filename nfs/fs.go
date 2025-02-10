package nfs

import (
	"errors"
	"github.com/allape/gohtvfs"
	"github.com/go-git/go-billy/v5"
	"io"
	"net/url"
	"os"
	"strings"
)

func NewBillyDufs(dufs *gohtvfs.DufsVFS) billy.Filesystem {
	return &BillyDufs{
		dufs: dufs,
	}
}

var NotImplError = errors.New("not implemented")

type BillyDufs struct {
	billy.Filesystem
	dufs *gohtvfs.DufsVFS
}

// region Basic

func (d BillyDufs) Create(filename string) (billy.File, error) {
	return d.Open(filename)
}

func (d BillyDufs) Open(filename string) (billy.File, error) {
	file, err := d.dufs.Open(filename)
	return &BillyDufsFile{file: file.(*gohtvfs.DufsFile)}, err
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

func (d BillyDufs) Lstat(filename string) (os.FileInfo, error) {
	return d.dufs.Stat(filename)
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
	file *gohtvfs.DufsFile
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

func (f *BillyDufsFile) Write(p []byte) (n int, err error) {
	return f.file.Write(p)
}

func (f *BillyDufsFile) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

func (f *BillyDufsFile) ReadAt(p []byte, off int64) (n int, err error) {
	_, err = f.file.Seek(off, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return f.file.Read(p)
}

func (f *BillyDufsFile) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

func (f *BillyDufsFile) Close() error {
	return f.file.Close()
}
