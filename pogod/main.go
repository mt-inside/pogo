package main

//go:generate protoc -I ../proto --go_out=plugins=grpc:../proto ../proto/pogo.proto

import (
	"log"
	"net"

	"github.com/mt-inside/pogo/pogod/cmd"

	pb "github.com/mt-inside/pogo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TODO: split repos (can inclue from another repo easily in go), vendor
// protos
const (
	port string = ":50001"
)

func main() {
	sock, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	pb.RegisterPogoServer(srv, &cmd.PogoServer{})
	pb.RegisterTasksServer(srv, &cmd.TasksServer{})
	// Turn on reflection so that clients can dynamically query our services
	reflection.Register(srv)
	log.Printf("serving on %v", port)
	if err := srv.Serve(sock); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
