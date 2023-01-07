/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package tree

import (
	"fmt"
	"io"

	"github.com/mandelsoft/filepath/pkg/filepath"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
)

type TreeWriter interface {
	Document(refpath string) (io.WriteCloser, string, error)
	Close() error
}

type fileTreeWriter struct {
	root string
	fs   vfs.FileSystem
}

func NewFileTreeWriter(path string, fss ...vfs.FileSystem) (TreeWriter, error) {
	fs := osfs.New()
	for _, f := range fss {
		if f != nil {
			fs = f
			break
		}
	}

	err := fs.MkdirAll(path, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot target dir %s: %w", path, err)
	}
	return &fileTreeWriter{
		root: path,
		fs:   fs,
	}, nil
}

func (w *fileTreeWriter) Document(refpath string) (io.WriteCloser, string, error) {
	path := w.root + refpath + ".md"
	err := w.fs.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return nil, path, fmt.Errorf("cannot create dir %s: %w", filepath.Dir(path), err)
	}
	f, err := w.fs.OpenFile(path, vfs.O_WRONLY|vfs.O_TRUNC|vfs.O_CREATE, 0644)
	return f, path, err
}

func (w *fileTreeWriter) Close() error {
	return nil
}
