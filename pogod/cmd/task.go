package cmd

/* Marshalls and unmarshalls types from gRPC to internal types (mostly
* primitives) */

import (
	"context"
	"log"

	"github.com/mt-inside/pogo/pogod/tasks"

	"github.com/mt-inside/pogo/pogod/model"
	"github.com/mt-inside/pogo/pogod/task"
	pb "github.com/mt-inside/pogo/proto"
)

type TasksServer struct{}

func (s *TasksServer) Add(ctxt context.Context, t *pb.Task) (*pb.Unit, error) {
	log.Printf("Adding task %v", t)
	new_id := model.NextTaskId()
	task := task.NewTask(new_id, t.Title, t.Category) /* Ignore user-specified ID, if it were even set */
	tasks.Add(task)

	return &pb.Unit{}, nil
}

func (s *TasksServer) List(f *pb.TaskFilter, stream pb.Tasks_ListServer) error {
	log.Println("Listing tasks")
	if (f.Fields &^ pb.TaskFields_category &^ pb.TaskFields_state &^ pb.TaskFields_type) != 0 {
		panic("unsupported filter option")
	}

	for _, t := range tasks.List() {
		include := true

		if f.Fields&pb.TaskFields_id != 0 {
			include = false
		}
		if f.Fields&pb.TaskFields_title != 0 {
			include = false
		}
		if f.Fields&pb.TaskFields_category != 0 &&
			f.Task.Category != t.Category {
			include = false
		}
		if f.Fields&pb.TaskFields_state != 0 &&
			int32(f.Task.State)&int32(t.State) == 0 {

			include = false
		}
		if f.Fields&pb.TaskFields_type != 0 &&
			int32(f.Task.Type) != int32(t.Type) {
			include = false
		}

		if include {
			if err := stream.Send(t.ToPB()); err != nil {
				return err
			}
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
