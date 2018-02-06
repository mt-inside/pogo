package tasks

import (
	"context"
	"io"
	"log"

	pb "github.com/mt-inside/pogo/proto"
)

var (
	pogo pb.PogoClient
)

/* This should convert PBs to and from internal types (so they can be
* rendered by the effectful layer), but since we don't
* have any yet, just deal externally in PBs */

func NewTaskClientHack(p pb.PogoClient) {
	pogo = p
}

func State() *pb.PogoState {
	state, err := pogo.GetState(context.Background(), &pb.Unit{})
	if err != nil {
		log.Fatalf("%v.State(_) = _, %v", pogo, err)
	}
	return state
}

func AddTask(t *pb.ProtoTask) {
	_, err := pogo.Add(context.Background(), t)
	if err != nil {
		log.Fatalf("%v.List(_) = _, %v", pogo, err)
	}
}

func ListTasks() (ts []*pb.Task) {
	// Background == default (no cancel, timeout, etc)
	stream, err := pogo.List(context.Background(), &pb.Unit{})
	if err != nil {
		log.Fatalf("%v.List(_) = _, %v", pogo, err)
	}
	for {
		t, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.List(_) = _, %v", pogo, err)
		}
		ts = append(ts, t)
	}
	return
}

func StartTask(id int64) {
	_, err := pogo.Start(context.Background(), &pb.Id{id})
	if err != nil {
		log.Fatalf("%v.Start(_) = _, %v", pogo, err)
	}
}

func CompleteTask(id int64) {
	_, err := pogo.Complete(context.Background(), &pb.Id{id})
	if err != nil {
		log.Fatalf("%v.Complete(_) = _, %v", pogo, err)
	}
}
