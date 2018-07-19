// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The nodefs package offers a high level API that resembles the
// kernel's idea of what an FS looks like.  File systems can have
// multiple hard-links to one file, for example. It is also suited if
// the data to represent fits in memory: you can construct the
// complete file system tree at mount time
package server

import (
	"time"

	"github.com/hanwen/go-fuse/fuse"
)

// A File object is returned from FileSystem.Open and
// FileSystem.Create.  Include the NewDefaultFile return value into
// the struct to inherit a null implementation.
type File interface {
	Read(off int64) ([]byte, fuse.Status)
}

// Wrap a File return in this to set FUSE flags.  Also used internally
// to store open file data.
type WithFlags struct {
	File

	// For debugging.
	Description string

	// Put FOPEN_* flags here.
	FuseFlags uint32

	// O_RDWR, O_TRUNCATE, etc.
	OpenFlags uint32
}

// Options contains time out options for a node FileSystem.  The
// default copied from libfuse and set in NewMountOptions() is
// (1s,1s,0s).
type Options struct {
	EntryTimeout    time.Duration
	AttrTimeout     time.Duration
	NegativeTimeout time.Duration

	// If set, replace all uids with given UID.
	// NewOptions() will set this to the daemon's
	// uid/gid.
	*fuse.Owner

	// This option exists for compatibility and is ignored.
	PortableInodes bool

	// If set, print debug information.
	Debug bool
}
