// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"github.com/LK4D4/grfuse/pb"
	"github.com/hanwen/go-fuse/fuse"
	"io/ioutil"
	"os"
	"path/filepath"
)

type loopbackFileSystem struct {
	// TODO - this should need default fill in.
	FileSystem
	Root string
}

// A FUSE filesystem that shunts all request to an underlying file
// system.  Its main purpose is to provide test coverage without
// having to build a synthetic filesystem.
func NewLoopbackFileSystem(root string) FileSystem {
	// Make sure the Root path is absolute to avoid problems when the
	// application changes working directory.
	root, err := filepath.Abs(root)
	if err != nil {
		panic(err)
	}
	return &loopbackFileSystem{
		FileSystem: NewDefaultFileSystem(),
		Root:       root,
	}
}

func (fs *loopbackFileSystem) GetPath(relPath string) string {
	return filepath.Join(fs.Root, relPath)
}

func (fs *loopbackFileSystem) GetAttr(name string, context *fuse.Context) (a *pb.Attr, code fuse.Status) {
	stat, err := os.Stat(fs.GetPath(name))
	if err != nil {
		return nil, fuse.ENOSYS
	}
	mode := uint32(stat.Mode().Perm())
	if stat.IsDir() {
		mode |= 0x4000
	} else {
		mode |= 0x8000
	}
	ret := pb.Attr{
		Mode:     mode,
		SizeAttr: uint64(stat.Size()),
	}
	return &ret, fuse.OK
}

func (fs *loopbackFileSystem) OpenDir(name string, context *fuse.Context) ([]pb.DirEntry, fuse.Status) {
	// What other ways beyond O_RDONLY are there to open
	// directories?
	files, err := ioutil.ReadDir(fs.GetPath(name))
	if err != nil {
		return nil, fuse.ENOSYS
	}

	var stream []pb.DirEntry
	for _, file := range files {
		stream = append(stream, pb.DirEntry{
			Name: file.Name(),
			Mode: uint32(file.Mode()),
		})
	}

	return stream, fuse.OK
}

func (fs *loopbackFileSystem) Open(name string, flags uint32, context *fuse.Context) (fuseFile File, status fuse.Status) {
	fd, err := os.Open(fs.GetPath(name))
	if err != nil {
		return nil, fuse.ENOSYS
	}
	bin, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, fuse.ENOSYS
	}
	return &foobar{bin}, fuse.OK
}

type foobar struct {
	payload []byte
}

func (f *foobar) Read(off int64) ([]byte, fuse.Status) {
	return f.payload, fuse.OK
}
