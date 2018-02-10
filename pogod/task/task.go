package task

import (
	pb "github.com/mt-inside/pogo/proto"
)

type TaskState int

const (
	Todo TaskState = iota
	Done
)

type Task struct {
	Id    int64
	Title string
	State TaskState
}

func (t *Task) ToPB() *pb.Task {
	return &pb.Task{Id: &pb.Id{t.Id}, Title: t.Title, State: pb.TaskState(t.State)}
}

func NewTask(idx int64, title string) *Task {
	return &Task{Id: idx, Title: title}
}
