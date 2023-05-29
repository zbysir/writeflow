package easyfs

import (
	stdFs "io/fs"
	"io/ioutil"
	"path"
	"path/filepath"
	"sort"
)

// 提供最简单的文件操作 API

type File struct {
	Name      string `json:"name"`
	Path      string `json:"path"` // full path
	DirPath   string `json:"dir_path"`
	IsDir     bool   `json:"is_dir"`
	CreatedAt int64  `json:"created_at"`
	ModifyAt  int64  `json:"modify_at"`
	Body      string `json:"body"`
}

type FileTree struct {
	File
	Items []FileTree `json:"items"`
}

func GetFile(fs stdFs.FS, path string) (fi *File, err error) {
	nf, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer nf.Close()
	finfo, err := nf.Stat()
	if err != nil {
		return nil, err
	}

	bs, err := ioutil.ReadAll(nf)
	if err != nil {
		return
	}

	dir, name := filepath.Split(path)
	return &File{
		Name:      name,
		Path:      path,
		DirPath:   dir,
		IsDir:     finfo.IsDir(),
		CreatedAt: 0,
		ModifyAt:  finfo.ModTime().Unix(),
		Body:      string(bs),
	}, nil
}

func GetFileTree(fs stdFs.FS, base string, deep int) (ft FileTree, err error) {
	_, ft.Name = path.Split(base)
	ft.Path = base
	ft.IsDir = true
	ft.DirPath = base
	if deep == 0 {
		return
	}

	fds, err := stdFs.ReadDir(fs, base)
	if err != nil {
		return ft, err
	}
	// 排序 文件夹在前
	sort.Slice(fds, func(i, j int) bool {
		iIsDir := fds[i].IsDir()
		jIsDir := fds[j].IsDir()
		if iIsDir && !jIsDir {
			return true
		} else if !iIsDir && jIsDir {
			return false
		}

		return fds[i].Name() < fds[j].Name()
	})
	for _, fd := range fds {
		if fd.Name() == ".git" {
			continue
		}

		srcfp := path.Join(base, fd.Name())

		if fd.IsDir() {
			ftw, err := GetFileTree(fs, srcfp, deep-1)
			if err != nil {
				return ft, err
			}
			ft.Items = append(ft.Items, ftw)
		} else {
			info, err := fd.Info()
			if err != nil {
				return ft, err
			}
			ft.Items = append(ft.Items, FileTree{
				File: File{
					Name:      fd.Name(),
					Path:      srcfp,
					DirPath:   base,
					IsDir:     false,
					CreatedAt: 0,
					ModifyAt:  info.ModTime().Unix(),
					Body:      "",
				},
				Items: nil,
			})
		}
	}

	return ft, nil
}
