package task

import (
	pb "github.com/mt-inside/pogo/api"
)

type TaskState int

const (
	Dummy1 TaskState = 0
	TODO             = 1
	DONE             = 2
)

type TaskType int

const (
	Dummy2 TaskType = 0
	TASK            = 1
	BREAK           = 2
)

type Task struct {
	Id       int64
	Title    string
	Category string
	State    TaskState
	Type     TaskType
}

func (t *Task) ToPB() *pb.Task {
	return &pb.Task{
		Id:       &pb.Id{Idx: t.Id},
		Title:    t.Title,
		Category: t.Category,
		State:    pb.TaskState(t.State),
		Type:     pb.TaskType(t.Type),
	}
}

func NewTask(idx int64, title string, category string) *Task {
	return &Task{
		Id:       idx,
		Title:    title,
		Category: category,
		State:    TODO,
		Type:     TASK,
	}
}
