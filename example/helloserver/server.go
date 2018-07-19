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
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: %s <mountpath>")
		os.Exit(1)
	}
	root := os.Args[1]

	loop := server.NewLoopbackFileSystem("C:\\Users\\admin")
	//hfs := &HelloFs{FileSystem: pathfs.NewDefaultFileSystem()}

	l, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterPathFSServer(s, server.New(loop))
	go s.Serve(l)
	log.Printf("Listen on %s for dir %s", l.Addr(), root)
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	for range sigCh {
		s.Stop()
		os.Exit(0)
	}
}
