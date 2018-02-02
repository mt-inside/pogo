package main

//go:generate protoc -I ../proto --go_out=plugins=grpc:../proto ../proto/pogo.proto

import (
	"log"
	"net"

	pb "github.com/mt-inside/pogo/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port string = ":50001"
)

type pogo_server struct{}

func (s *pogo_server) Add(ctxt context.Context, t *pb.Task) (*pb.Unit, error) {
	return &pb.Unit{}, nil
}

func (s *pogo_server) List(_ *pb.Unit, stream pb.Pogo_ListServer) (error) {
	if err := stream.Send(&pb.Task{"get up!"}); err != nil {
		return err
	}
	return nil
}

func main() {
	sock, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	pb.RegisterPogoServer(srv, &pogo_server{})
	reflection.Register(srv)
	log.Printf("serving on %v", port)
	if err := srv.Serve(sock); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
