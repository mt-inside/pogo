package cmd

import (
	"context"
	"log"

	"github.com/mt-inside/pogo/pogod/tasks"

	pb "github.com/mt-inside/pogo/proto"
)

type PogoServer struct{}

// This is really looking like it should be an MVC...
// TODO: unmarshall grpc; convert to internal types

func (s *PogoServer) GetState(ctxt context.Context, _ *pb.Unit) (*pb.PogoState, error) {
	state, task, remain := tasks.State()

	var ret *pb.PogoState
	if state == tasks.Idle {
		ret = &pb.PogoState{
			State: pb.PogoState_State(state),
		}
	} else {
		ret = &pb.PogoState{
			State:         pb.PogoState_State(state),
			Task:          task.ToPB(),
			RemainingTime: remain,
		}
	}

	return ret, nil
}

func (s *PogoServer) Add(ctxt context.Context, pt *pb.ProtoTask) (*pb.Unit, error) {
	log.Printf("Adding task %v", pt)
	tasks.Add(pt.Title)

	return &pb.Unit{}, nil
}

func (s *PogoServer) List(_ *pb.Unit, stream pb.Pogo_ListServer) error {
	log.Println("Listing tasks")
	for _, t := range tasks.List() {
		if err := stream.Send(t.ToPB()); err != nil {
			return err
		}
	}
	return nil
}

func (s *PogoServer) Start(ctxt context.Context, id *pb.Id) (*pb.Unit, error) {
	log.Printf("Request to start task %v", id)

	tasks.Start(id.Idx)

	return &pb.Unit{}, nil
}

func (s *PogoServer) Complete(ctxt context.Context, id *pb.Id) (*pb.Unit, error) {
	log.Printf("Request to complete task %v", id)

	tasks.Complete(id.Idx)

	return &pb.Unit{}, nil
}
