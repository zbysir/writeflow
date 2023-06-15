package fshook

import (
	"io"
	"io/fs"
)

func (f *FsHook) Open(name string) (fs.File, error) {
	i, err := f.i.Open(name)
	if err != nil {
		return nil, err
	}

	c, ok := f.inject[name]
	if ok {
		return NewFileHook(i, c), nil
	}

	return i, err
}

type FsHook struct {
	i      fs.FS
	inject map[string]func(body []byte) []byte
}

func NewFsHook(i fs.FS, inject map[string]func(body []byte) []byte) *FsHook {
	if inject == nil {
		inject = make(map[string]func(body []byte) []byte)
	}
	return &FsHook{i: i, inject: inject}
}

type FileInfoHook struct {
	fs.FileInfo
	size int
}

func (f *FileInfoHook) Size() int64 {
	return int64(f.size)
}

func (f *FileHook) Stat() (fs.FileInfo, error) {
	info, err := f.i.Stat()
	if err != nil {
		return nil, err
	}
	return &FileInfoHook{FileInfo: info, size: len(f.buf)}, nil
}

func (f *FileHook) Read(ibytes []byte) (int, error) {
	copy(ibytes, f.buf)
	l := 0
	if len(f.buf) > len(ibytes) {
		f.buf = f.buf[len(ibytes):]
		l = len(ibytes)
	} else {
		l = len(f.buf)
		f.buf = nil
	}

	return l, nil
}

func (f *FileHook) Close() error {
	return f.i.Close()
}

type FileHook struct {
	i   fs.File
	buf []byte
}

func NewFileHook(i fs.File, inject func(body []byte) []byte) *FileHook {
	bs, err := io.ReadAll(i)
	if err != nil {
		return &FileHook{i: i}
	}

	bs = inject(bs)

	return &FileHook{i: i, buf: bs}
}
