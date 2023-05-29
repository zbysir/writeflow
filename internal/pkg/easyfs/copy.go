package easyfs

import (
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"io"
	"io/fs"
	"os"
	"path"
)

func CopyFile(src, dst string, srcFs fs.FS, dstFs billy.Filesystem) error {
	var err error
	var srcfd fs.File
	var dstfd billy.File

	if srcfd, err = srcFs.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = dstFs.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}

	return nil
}

func CopyDir(src string, dst string, srcFs fs.FS, dstFs billy.Filesystem) error {
	var err error
	var fds []os.DirEntry
	var srcinfo os.FileInfo

	if srcinfo, err = fs.Stat(srcFs, src); err != nil {
		return fmt.Errorf("fs.State %s error: %w", src, err)
	}
	if !srcinfo.IsDir() {
		return CopyFile(src, dst, srcFs, dstFs)
	}

	if err = dstFs.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}
	if fds, err = fs.ReadDir(srcFs, src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp, srcFs, dstFs); err != nil {
				log.Warnf("exportDir error: %s", err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp, srcFs, dstFs); err != nil {
				log.Warnf("exportFile error: %s", err)
			}
		}
	}
	return nil
}
