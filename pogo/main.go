package main

//go:generate protoc -I ../proto --go_out=plugins=grpc:../proto ../proto/pogo.proto

import (
	"context"
	"io"
	"log"

	pb "github.com/mt-inside/pogo/proto"
	"google.golang.org/grpc"
)

const (
	serverAddr string = "localhost:50001"
)

func addTask(pogo pb.PogoClient, t *pb.Task) {
	_, err := pogo.Add(context.Background(), t)
	if err != nil {
		log.Fatalf("%v.List(_) = _, %v", pogo, err)
	}
}

func listTasks(pogo pb.PogoClient) {
	// Background == default (no cancel, timeout, etc)
	stream, err := pogo.List(context.Background(), &pb.Unit{})
	if err != nil {
		log.Fatalf("%v.List(_) = _, %v", pogo, err)
	}
	for {
		task, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.List(_) = _, %v", pogo, err)
		}
		log.Println(task)
	}
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("Couldn't connect to pogo server: %v", err)
	}
	defer conn.Close()
	pogo := pb.NewPogoClient(conn)

	addTask(pogo, &pb.Task{"get up"})
	addTask(pogo, &pb.Task{"go climbing"})
	listTasks(pogo)
}
