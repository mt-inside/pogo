package main

//go:generate protoc -I ../proto --go_out=plugins=grpc:../proto ../proto/pogo.proto

import (
	"log"

	"github.com/mt-inside/pogo/pogo/cmd"
	"github.com/mt-inside/pogo/pogo/tasks"

	pb "github.com/mt-inside/pogo/proto"
	"google.golang.org/grpc"
)

const (
	serverAddr string = "localhost:50001"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("Couldn't connect to pogo server: %v", err)
	}
	defer conn.Close()
	// TODO: nope
	pogo := pb.NewPogoClient(conn)
	tasks.NewTaskClientHack(pogo)
	cmd.Execute()
}
