package main

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/samber/mo"
)

type Directory struct {
	path string
	fsys fs.FS
}

type Order int

const (
	ORDER_NAME Order = iota
	ORDER_TIME
)

func NewDirectory(path string) Directory {
	return Directory{path: path, fsys: os.DirFS(path)}
}

func (d Directory) String() string {
	return d.path
}

func (d Directory) Name() string {
	return filepath.Base(d.String())
}

func (d Directory) IsHidden() bool {
	return strings.HasPrefix(d.Name(), ".")
}

func (d Directory) Parent() mo.Option[Directory] {
	parent := path.Dir(d.String())

	if parent == d.String() {
		return mo.None[Directory]()
	} else {
		return mo.Some(NewDirectory(parent))
	}

}

func (d Directory) Dirs(showAll bool, order Order) mo.Result[[]Directory] {
	files, err := fs.ReadDir(d.fsys, ".")

	if err != nil {
		return mo.Err[[]Directory](err)
	}

	var entries []fs.DirEntry
	for _, file := range files {
		if file.IsDir() || getSymlink(filepath.Join(d.String(), file.Name())).IsPresent() {
			if showAll || !strings.HasPrefix(file.Name(), ".") {
				entries = append(entries, file)
			}
		}
	}

	switch order {
	case ORDER_NAME:
	case ORDER_TIME:
		sort.Slice(entries, func(i, j int) bool {
			af, _ := entries[i].Info()
			bf, _ := entries[j].Info()
			return bf.ModTime().After(af.ModTime())
		})
	}

	var directories []Directory
	for _, entry := range entries {
		directories = append(directories, NewDirectory(filepath.Join(d.String(), entry.Name())))
	}

	return mo.Ok(directories)
}

func (d Directory) SymLink() mo.Option[string] {
	return getSymlink(d.String())
}

func getSymlink(path string) mo.Option[string] {
	info, err := os.Lstat(path)

	if err != nil {
		return mo.None[string]()
	}

	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		realPath, err := os.Readlink(path)
		if err != nil {
			return mo.None[string]()
		}

		info, err := os.Lstat(realPath)

		if err != nil {
			return mo.None[string]()
		}

		if info.IsDir() {
			return mo.Some(realPath)
		} else {
			return mo.None[string]()
		}

	} else {
		return mo.None[string]()
	}
}
