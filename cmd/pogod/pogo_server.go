package main

import (
	"context"
	"log"

	"github.com/mt-inside/pogo/pkg/pogod/tasks"

	pb "github.com/mt-inside/pogo/api"
)

type PogoServer struct{}

func (s *PogoServer) GetStatus(ctxt context.Context, _ *pb.Unit) (*pb.Status, error) {
	log.Println("Getting system status")
	state, task, remain := tasks.Status()

	var ret *pb.Status
	if state == tasks.Idle {
		ret = &pb.Status{
			State: pb.Status_SystemState(state),
		}
	} else {
		ret = &pb.Status{
			State:         pb.Status_SystemState(state),
			Task:          task.ToPB(),
			RemainingTime: remain,
		}
	}

	return ret, nil
}
