package cmd

/* Marshalls and unmarshalls types from gRPC to internal types (mostly
* primitives) */

import (
	"context"
	"log"

	"github.com/mt-inside/pogo/pogod/tasks"

	pb "github.com/mt-inside/pogo/proto"
)

type TasksServer struct{}

func (s *TasksServer) Add(ctxt context.Context, t *pb.Task) (*pb.Unit, error) {
	log.Printf("Adding task %v", t)
	tasks.Add(t.Title) /* Ignore ID, if it were even set */

	return &pb.Unit{}, nil
}

func (s *TasksServer) List(f *pb.TaskFilter, stream pb.Tasks_ListServer) error {
	log.Println("Listing tasks")
	if f.Fields != 0 {
		panic("unsupported filter option")
	}
	for _, t := range tasks.List() {
		if err := stream.Send(t.ToPB()); err != nil {
			return err
		}
	}
	return nil
}

func (s *TasksServer) Start(ctxt context.Context, id *pb.Id) (*pb.Unit, error) {
	log.Printf("Request to start task %v", id)

	if err := tasks.Start(id.Idx); err != nil {
		return nil, err
	} else {
		return &pb.Unit{}, nil
	}
}

func (s *TasksServer) Stop(ctxt context.Context, _ *pb.Unit) (*pb.Unit, error) {
	log.Printf("Request to stop task")

	if err := tasks.Stop(); err != nil {
		return nil, err
	} else {
		return &pb.Unit{}, nil
	}

	return &pb.Unit{}, nil
}

func (s *TasksServer) Complete(ctxt context.Context, id *pb.Id) (*pb.Unit, error) {
	log.Printf("Request to complete task %v", id)

	if err := tasks.Complete(id.Idx); err != nil {
		return nil, err
	} else {
		return &pb.Unit{}, nil
	}

	return &pb.Unit{}, nil
}
