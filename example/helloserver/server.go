package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"github.com/LK4D4/grfuse/pb"
	"github.com/LK4D4/grfuse/server"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

type HelloFs struct {
	pathfs.FileSystem
}

func (me *HelloFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	switch name {
	case "file.txt":
		return &fuse.Attr{
			Mode: fuse.S_IFREG | 0644, Size: uint64(len(name)),
		}, fuse.OK
	case "":
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (me *HelloFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	if name == "" {
		c = []fuse.DirEntry{{Name: "file.txt", Mode: fuse.S_IFREG}}
		return c, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (me *HelloFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	if name != "file.txt" {
		return nil, fuse.ENOENT
	}
	if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}
	return nodefs.NewDataFile([]byte(name)), fuse.OK
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: %s <mountpath>")
		os.Exit(1)
	}
	root := os.Args[1]

	hfs := &HelloFs{FileSystem: pathfs.NewDefaultFileSystem()}

	nfs := pathfs.NewPathNodeFs(hfs, nil)
	fuseSrv, _, err := nodefs.MountRoot(root, nfs.Root(), nil)
	if err != nil {
		log.Fatal(err)
	}
	go fuseSrv.Serve()
	l, err := net.Listen("tcp", "127.0.0.1:50000")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterPathFSServer(s, server.New(hfs))
	go s.Serve(l)
	log.Printf("Listen on %s for dir %s", l.Addr(), root)
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	for range sigCh {
		fuseSrv.Unmount()
		s.Stop()
		os.Exit(0)
	}
}
