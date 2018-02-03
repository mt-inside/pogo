package task

import (
	"context"
	"io"
	"log"

	pb "github.com/mt-inside/pogo/proto"
)

var (
	pogo pb.PogoClient
)

func NewTaskClientHack(p pb.PogoClient) {
	pogo = p
}

func AddTask(t *pb.Task) {
	_, err := pogo.Add(context.Background(), t)
	if err != nil {
		log.Fatalf("%v.List(_) = _, %v", pogo, err)
	}
}

func ListTasks() {
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
