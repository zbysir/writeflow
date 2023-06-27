package gobilly

import (
	"bytes"
	"fmt"
	"github.com/docker/libkv/store"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/helper/chroot"
	"github.com/thoas/go-funk"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// 考虑使用 gobilly 作为文件系统抽象，fusefs 有点难用。
// https://github.com/go-git/go-billy

var _ billy.Filesystem = (*DbFs)(nil)

const (
	defaultDirectoryMode = 0755
	defaultCreateMode    = 0666
)

type DbFs struct {
	store store.Store
	root  string
}

func NewDbFs(store store.Store) *DbFs {
	d := &DbFs{store: store, root: "/"}
	err := d.MkdirAll("", 0)
	if err != nil {
		log.Errorf("mkdir root '/' error: %v", err)
	}
	return d
}

func (d *DbFs) Create(filename string) (billy.File, error) {
	// mkdir
	dir := filepath.Dir(filename)
	err := d.MkdirAll(dir, 0)
	if err != nil {
		return nil, err
	}
	f := NewFile(d.store, filepath.Join(d.root, filename))
	_, err = f.Write(nil)
	if err != nil {
		return nil, err
	}
	return f, err
}

func (d *DbFs) Open(filename string) (billy.File, error) {
	filename = filepath.Join(d.root, filename)
	f := NewFile(d.store, filename)
	err := f.loadOnce()
	if err != nil {
		return nil, err
	}
	log.Debugf("open file or dir: %v %v", filename, f.Mode().IsDir())

	return f, nil
}

func (d *DbFs) OpenFile(filename string, flag int, perm os.FileMode) (billy.File, error) {
	return d.Open(filename)
}

// Stat TODO 将 FileInfo 与文件分开存储，避免每次读取信息都读大文件内容
func (d *DbFs) Stat(filename string) (os.FileInfo, error) {
	filename = filepath.Join(d.root, filename)
	f := NewFile(d.store, filename)
	err := f.loadOnce()
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (d *DbFs) Rename(oldName, newName string) error {
	oldName = path.Join(d.root, oldName)
	newName = path.Join(d.root, newName)
	kv, err := d.store.Get(oldName)
	if err != nil {
		return err
	}

	var undo func()
	newKv, err := d.store.Get(newName)
	if err == nil {
		undo = func() {
			d.store.Put(oldName, kv.Value, nil)
			d.store.Put(newKv.Key, newKv.Value, nil)
		}
	} else {
		undo = func() {
			d.store.Put(oldName, kv.Value, nil)
			d.store.Delete(newName)
		}
	}

	if err := d.store.Put(newName, kv.Value, nil); err != nil {
		log.Errorf("%v", err)
		return err
	}

	if err := d.store.Delete(oldName); err != nil {
		undo()
		return err
	}

	return nil
}

// Remove 支持删除 文件 和 文件夹
// TODO 测试删除文件夹
func (d *DbFs) Remove(filename string) error {
	filename = path.Join(d.root, filename)

	if err := d.store.Delete(filename); err != nil {
		return err
	}

	//if err := fs.kvStore.DeleteTree(name); err != nil {
	//	logrus.Error(err)
	//	return fuse.EIO
	//}
	return nil
}

func (d *DbFs) Join(elem ...string) string {
	return path.Join(elem...)
}

func (d *DbFs) TempFile(dir, prefix string) (billy.File, error) {
	if dir == "" {
		dir = "temp"
	}
	fname := path.Join(dir, prefix, funk.RandomString(12, []rune("abcdefghijklmnopqrstuvwxyz0123456789")))
	return d.Open(fname)
}

func (d *DbFs) ReadDir(path string) ([]os.FileInfo, error) {
	path = filepath.Join(d.root, path)
	log.Debugf("readdir %+v", path)
	// get dir file
	_, err := d.store.Get(path)
	if err != nil {
		if err == store.ErrKeyNotFound {
			return nil, os.ErrNotExist
		}
		return nil, err
	}

	kvs, err := d.store.List(path)
	if err != nil {
		return nil, fmt.Errorf("store.list error: %w", err)
	}

	var entries []os.FileInfo
	for _, kv := range kvs {
		fullPath := kv.Key
		if path != "/" {
			if !strings.HasPrefix(fullPath, path+"/") {
				continue
			}
		}

		dir, fi := filepath.Split(fullPath)
		if dir == path && fi == "" {
			log.Debugf("skipping base %s", fullPath)
			continue
		}

		// 相对于入参的路径
		relatPath := strings.TrimPrefix(fullPath, path)
		relatPath = strings.TrimPrefix(relatPath, "/")

		// 不允许有 /，除了只允许最后一个
		// /src
		// /src/jsx ok
		// /src/js/ ok
		// /src/js/x fail

		index := strings.Index(relatPath, "/")
		if index != -1 && index != len(relatPath)-1 {
			log.Debugf("skipping subtree %s %s", dir, fullPath)
			continue
		}

		rel, err := filepath.Rel(d.root, fullPath)
		if err != nil {
			return nil, err

		}
		log.Debugf("rel %v", rel)
		finfo, err := d.Stat(rel)
		if err != nil {
			return nil, fmt.Errorf("stat error: %w [%v]", err, rel)
		}

		entries = append(entries, finfo)
	}

	// 排序 文件夹在前
	sort.Slice(entries, func(i, j int) bool {
		a := entries[i]
		b := entries[j]
		iIsDir := a.Mode().IsDir()
		jIsDir := b.Mode().IsDir()
		if iIsDir && !jIsDir {
			return true
		} else if !iIsDir && jIsDir {
			return false
		}

		return a.Name() < b.Name()
	})
	return entries, nil
}

// MkdirAll 传递 ” 和 '/' 都将会创建 root 目录
func (d *DbFs) MkdirAll(filename string, perm os.FileMode) error {
	filename = filepath.Clean(filename)
	if filename == "." {
		filename = ""
	}
	if filename == "/" {
		filename = ""
	}
	log.Debugf("MkdirAll '%v'", filename)

	// 先检查是否存在，如果存在则不处理
	finfo, err := d.Stat(filename)
	if err != nil {
		// 不存在则新建
		if err == os.ErrNotExist {
			if filename != "" {
				dir := filepath.Dir(filename)
				// 新建上级
				err = d.MkdirAll(dir, perm)
				if err != nil {
					return err
				}
			}

			// 新建本级
			f := NewFile(d.store, filepath.Join(d.root, filename))
			err = f.WriteDir()
			if err != nil {
				return err
			}

		} else {
			return err
		}
	} else {
		if !finfo.IsDir() {
			return fmt.Errorf("file '%v' exist", path.Join(d.root, filename))
		}
	}
	log.Debugf("MkdirAll '%v' done", filename)

	return nil
}

func (d *DbFs) Lstat(filename string) (os.FileInfo, error) {
	return d.Stat(filename)
}

// Symlink 暂时不实现
func (d *DbFs) Symlink(target, link string) error {
	return nil
}

func (d *DbFs) Readlink(link string) (string, error) {
	return "", nil
}

func (d *DbFs) Chroot(path string) (billy.Filesystem, error) {
	return chroot.New(d, path), nil
}

func (d *DbFs) Root() string {
	return d.root
}

type File struct {
	kvStore store.Store
	name    string
	content []byte
	attr    *FileAttr

	offset int
	// 一个文件只会读取一次，防止并发导致循环读的时候读到不同内容
	once      sync.Once
	loadError error
}

func (f *File) Type() fs.FileMode {
	return f.Mode()
}

func (f *File) Info() (fs.FileInfo, error) {
	return f, nil
}

var _ billy.File = (*File)(nil)
var _ os.FileInfo = (*File)(nil)

func NewFile(kvStore store.Store, name string) *File {
	return &File{kvStore: kvStore, name: name}
}

func (f *File) Size() int64 {
	err := f.loadOnce()
	if err != nil {
		return 0
	}
	return int64(f.attr.Size)
}

func (f *File) Mode() fs.FileMode {
	err := f.loadOnce()
	if err != nil {
		return 0
	}

	return f.attr.Mode
}

func (f *File) ModTime() time.Time {
	err := f.loadOnce()
	if err != nil {
		return time.Time{}
	}
	return time.Unix(int64(f.attr.MTime), 0)
}

func (f *File) IsDir() bool {
	return f.Mode().IsDir()
}

func (f *File) Sys() any {
	return nil
}

func (f *File) Name() string {
	_, name := filepath.Split(f.name)
	return name
}

func (f *File) loadOnce() error {
	f.once.Do(func() {
		kv, err := f.kvStore.Get(f.name)
		if err != nil {
			if err == store.ErrKeyNotFound {
				f.loadError = os.ErrNotExist
				var attr FileAttr
				attr.Mode = defaultCreateMode
				f.attr = &attr
			} else {
				f.loadError = err
			}
			return
		}

		value := kv.Value

		var attrx []byte
		var content = value

		// 第一行是元数据
		sp := bytes.SplitN(value, []byte("\n\n"), 2)
		if len(sp) == 2 {
			attrx = sp[0]
			content = sp[1]
		}
		var attr FileAttr

		if len(attrx) == 0 {
			// default
			attr.Mode = defaultCreateMode
		} else {
			attr.fromByte(attrx)
		}

		f.content = content
		f.attr = &attr
	})

	return f.loadError
}

func (f *File) Write(p []byte) (n int, err error) {
	err = f.loadOnce()
	if err != nil {
		if err == os.ErrNotExist {
			err = nil
		} else {
			return
		}
	} else {
		if f.IsDir() {
			return 0, fmt.Errorf("dir '%v' exist", f.name)
		}
	}

	a := f.attr
	if a.CTime == 0 {
		a.CTime = uint64(time.Now().Unix())
	}
	a.MTime = uint64(time.Now().Unix())
	a.Size = uint64(len(p))
	attrx := a.toByte()

	fileAll := append(attrx, '\n', '\n')
	fileAll = append(fileAll, p...)

	err = f.kvStore.Put(f.name, fileAll, nil)
	if err != nil {
		return
	}
	return len(p), nil
}

func (f *File) WriteDir() (err error) {
	err = f.loadOnce()
	if err != nil {
		if err == os.ErrNotExist {
			// 不存在则新建
			a := f.attr
			if a.CTime == 0 {
				a.CTime = uint64(time.Now().Unix())
			}
			a.Mode = defaultDirectoryMode | os.ModeDir
			a.MTime = uint64(time.Now().Unix())
			attrx := a.toByte()

			fileAll := append(attrx, '\n', '\n')

			err = f.kvStore.Put(f.name, fileAll, nil)
			if err != nil {
				return
			}
		} else {
			return
		}
	} else {
		if !f.attr.Mode.IsDir() {
			return fmt.Errorf("file '%v' exist", f.name)
		}
	}

	return nil
}

func (f *File) Read(buf []byte) (n int, err error) {
	err = f.loadOnce()
	if err != nil {
		return
	}

	if f.offset >= len(f.content) {
		return 0, io.EOF
	}

	end := f.offset + len(buf)
	if end > len(f.content) {
		end = len(f.content)
	}
	copy(buf, f.content[f.offset:end])

	n = end - f.offset
	f.offset += n
	return
}

func (f *File) ReadAt(buf []byte, off int64) (n int, err error) {
	err = f.loadOnce()
	if err != nil {
		return
	}
	if off >= int64(len(f.content)) {
		return 0, io.EOF
	}

	end := int(off) + len(buf)
	if end > len(f.content) {
		end = len(f.content)
	}
	copy(buf, f.content[off:end])

	return
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f *File) Close() error {
	return nil
}

func (f *File) Lock() error {
	return nil
}

func (f *File) Unlock() error {
	return nil
}

func (f *File) Truncate(size int64) error {
	return nil
}

type FileAttr struct {
	MTime uint64
	CTime uint64
	Size  uint64
	Mode  os.FileMode
}

func (a *FileAttr) toByte() []byte {
	return encode(a)
}

func (a *FileAttr) fromByte(bs []byte) {
	if len(bs) == 0 {
		return
	}
	decode(bs, a)
}
