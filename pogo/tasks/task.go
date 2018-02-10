package tasks

import (
	"context"
	"io"
	"log"

	"github.com/mt-inside/pogo/pogo/clients"

	pb "github.com/mt-inside/pogo/proto"
)

/* This should convert PBs to and from internal types (so they can be
* rendered by the effectful layer), but since we don't
* have any yet, just deal externally in PBs */

func GetStatus() *pb.Status {
	pogo := clients.GetPogoClient()
	status, err := pogo.GetStatus(context.Background(), &pb.Unit{})
	if err != nil {
		log.Fatalf("%v.State(_) = _, %v", pogo, err)
	}
	return status
}

func AddTask(t *pb.Task) {
	tasks := clients.GetTasksClient()
	_, err := tasks.Add(context.Background(), t)
	if err != nil {
		log.Fatalf("%v.List(_) = _, %v", tasks, err)
	}
}

func ListTasks() (ts []*pb.Task) {
	tasks := clients.GetTasksClient()
	// Background == default (no cancel, timeout, etc)
	stream, err := tasks.List(context.Background(), &pb.TaskFilter{})
	if err != nil {
		log.Fatalf("%v.List(_) = _, %v", tasks, err)
	}
	for {
		t, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.List(_) = _, %v", tasks, err)
		}
		ts = append(ts, t)
	}
	return
}

func StartTask(id int64) {
	tasks := clients.GetTasksClient()
	_, err := tasks.Start(context.Background(), &pb.Id{id})
	if err != nil {
		log.Fatalf("%v.Start(_) = _, %v", tasks, err)
	}
}

func StopTask() {
	tasks := clients.GetTasksClient()
	_, err := tasks.Stop(context.Background(), &pb.Unit{})
	if err != nil {
		log.Fatalf("%v.Stop(_) = _, %v", tasks, err)
	}
}

func CompleteTask(id int64) {
	tasks := clients.GetTasksClient()
	_, err := tasks.Complete(context.Background(), &pb.Id{id})
	if err != nil {
		log.Fatalf("%v.Complete(_) = _, %v", tasks, err)
	}
}
