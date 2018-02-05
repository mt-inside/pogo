package tasks

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
	return &pb.Task{Id: &pb.Id{t.Id}, Title: t.Title}
}
func NewTaskFromPBProto(pbt *pb.ProtoTask, idx int64) *Task {
	return &Task{Id: idx, Title: pbt.Title}
}
func NewTaskFromPB(pbt *pb.Task) *Task {
	return &Task{Id: pbt.Id.Idx, Title: pbt.Title}
}
