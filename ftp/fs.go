package ftp

import (
	"errors"
	vfs "github.com/allape/go-http-vfs"
	ftpserver "github.com/fclairamb/ftpserverlib"
	"github.com/spf13/afero"
	"os"
	"time"
)

var NotImplemented = errors.New("not implemented")

type DufsClientDriver struct {
	ftpserver.ClientDriver
	dufs *vfs.DufsVFS
}

func (d *DufsClientDriver) Create(name string) (afero.File, error) {
	return d.Open(name)
}

func (d *DufsClientDriver) Mkdir(name string, perm os.FileMode) error {
	return d.dufs.Mkdir(name, perm)
}

func (d *DufsClientDriver) MkdirAll(path string, perm os.FileMode) error {
	return d.dufs.Mkdir(path, perm)
}

func (d *DufsClientDriver) Open(name string) (afero.File, error) {
	file, err := d.dufs.Open(name)
	if err != nil {
		return nil, err
	}
	return &DufsAferoFile{
		file: file.(*vfs.DufsFile),
	}, nil
}

func (d *DufsClientDriver) OpenFile(name string, _ int, _ os.FileMode) (afero.File, error) {
	file, err := d.Open(name)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		return nil, os.ErrInvalid
	}

	return file, nil
}

func (d *DufsClientDriver) Remove(name string) error {
	return d.dufs.Remove(name)
}

func (d *DufsClientDriver) RemoveAll(path string) error {
	return d.dufs.Remove(path)
}

func (d *DufsClientDriver) Rename(oldname, newname string) error {
	return d.dufs.Rename(oldname, newname)
}

func (d *DufsClientDriver) Stat(name string) (os.FileInfo, error) {
	state, err := d.dufs.Stat(name)
	if err != nil {
		return nil, err
	}
	return &DufsAferoFileInfo{
		fileInfo: state,
	}, nil
}

func (d *DufsClientDriver) Name() string {
	return Name
}

func (d *DufsClientDriver) Chmod(_ string, _ os.FileMode) error {
	return NotImplemented
}

func (d *DufsClientDriver) Chown(_ string, _, _ int) error {
	return NotImplemented
}

func (d *DufsClientDriver) Chtimes(_ string, _ time.Time, _ time.Time) error {
	return NotImplemented
}

type DufsAferoFile struct {
	afero.File
	file *vfs.DufsFile
}

func (f *DufsAferoFile) Name() string {
	return f.file.Name
}

func (f *DufsAferoFile) Readdir(count int) ([]os.FileInfo, error) {
	stat, err := f.file.CachedStat()
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, os.ErrInvalid
	}

	files, err := f.file.ReadDir(count)
	if err != nil {
		return nil, err
	}

	fileInfos := make([]os.FileInfo, len(files))
	for i, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		fileInfos[i] = &DufsAferoFileInfo{
			fileInfo: info,
		}
	}

	return fileInfos, nil
}

func (f *DufsAferoFile) Readdirnames(n int) ([]string, error) {
	fileInfos, err := f.Readdir(n)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(fileInfos))
	for i, fileInfo := range fileInfos {
		names[i] = fileInfo.Name()
	}
	return names, nil
}

func (f *DufsAferoFile) Stat() (os.FileInfo, error) {
	stat, err := f.file.CachedStat()
	if err != nil {
		return nil, err
	}
	return &DufsAferoFileInfo{
		fileInfo: stat,
	}, nil
}

func (f *DufsAferoFile) Sync() error {
	_, err := f.file.Stat()
	return err
}

func (f *DufsAferoFile) Truncate(_ int64) error {
	return NotImplemented
}

func (f *DufsAferoFile) WriteString(s string) (int, error) {
	return f.file.Write([]byte(s))
}

func (f *DufsAferoFile) Close() error {
	return f.file.Close()
}

func (f *DufsAferoFile) Read(p []byte) (n int, err error) {
	stat, err := f.file.CachedStat()
	if err != nil {
		return 0, err
	} else if stat.IsDir() {
		return 0, errors.New("is a directory")
	}
	return f.file.Read(p)
}

func (f *DufsAferoFile) ReadAt(p []byte, off int64) (n int, err error) {
	return f.file.ReadAt(p, off)
}

func (f *DufsAferoFile) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

func (f *DufsAferoFile) Write(p []byte) (n int, err error) {
	return f.file.Write(p)
}

func (f *DufsAferoFile) WriteAt(p []byte, off int64) (n int, err error) {
	return f.file.WriteAt(p, off)
}

type DufsAferoFileInfo struct {
	os.FileInfo
	fileInfo os.FileInfo
}

func (i *DufsAferoFileInfo) Name() string {
	//if i.fileInfo.IsDir() && !strings.HasPrefix(i.fileInfo.Name(), "/") {
	//	return i.fileInfo.Name() + "/"
	//}
	return i.fileInfo.Name()
}

func (i *DufsAferoFileInfo) Size() int64 {
	if i.fileInfo.IsDir() {
		return 0
	}
	return i.fileInfo.Size()
}

func (i *DufsAferoFileInfo) Mode() os.FileMode {
	return i.fileInfo.Mode()
}

func (i *DufsAferoFileInfo) ModTime() time.Time {
	return i.fileInfo.ModTime()
}

func (i *DufsAferoFileInfo) IsDir() bool {
	return i.fileInfo.IsDir()
}

func (i *DufsAferoFileInfo) Sys() interface{} {
	return i.fileInfo.Sys()
}
