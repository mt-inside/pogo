package data

import (
	pb "github.com/mt-inside/pogo/proto"
)

var tasks []*pb.Task

func Add(t *pb.Task) {
	tasks = append(tasks, t)
}

func List() []*pb.Task {
	return tasks
}
