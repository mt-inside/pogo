package main

//go:generate protoc -I ../proto --go_out=plugins=grpc:../proto ../proto/pogo.proto

import (
	"context"
	"log"
	"net"

	"github.com/mt-inside/pogo/pogod/data"
	"github.com/mt-inside/pogo/pogod/tasks"

	pb "github.com/mt-inside/pogo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port string = ":50001"
)

type pogo_server struct{}

// This is really looking like it should be an MVC...

func (s *pogo_server) Add(ctxt context.Context, pt *pb.ProtoTask) (*pb.Unit, error) {
	log.Printf("Adding task %v", pt)
	data.Add(pt)

	return &pb.Unit{}, nil
}

func (s *pogo_server) List(_ *pb.Unit, stream pb.Pogo_ListServer) error {
	log.Println("Listing tasks")
	for _, t := range data.List() {
		if err := stream.Send(t.ToPB()); err != nil {
			return err
		}
	}
	return nil
}

func (s *pogo_server) Start(ctxt context.Context, id *pb.Id) (*pb.Unit, error) {
	log.Printf("Request to start task %v", id)
	t := data.Find(id.Idx)
	tasks.Start(t)
	log.Printf("Started task %v", t)

	return &pb.Unit{}, nil
}

func (s *pogo_server) Complete(ctxt context.Context, id *pb.Id) (*pb.Unit, error) {
	log.Printf("Request to complete task %v", id)
	t := data.Find(id.Idx)
	t.State = tasks.Done
	log.Printf("Completed task %v", t)

	return &pb.Unit{}, nil
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
