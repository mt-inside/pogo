package data

import (
	. "github.com/mt-inside/pogo/pogod/tasks"
	pb "github.com/mt-inside/pogo/proto"
)

var (
	tasks    map[int64]*Task = make(map[int64]*Task)
	next_idx int64           = 0
)

func Add(pt *pb.ProtoTask) {
	t := NewTaskFromPBProto(pt, next_idx)
	tasks[next_idx] = t

	next_idx += 1
}

func List() map[int64]*Task {
	return tasks
}

func Find(id int64) *Task {
	if t, found := tasks[id]; found {
		return t
	}
	panic("no such task id")
}

/*func List() (ts []*Task) {
    for _, t := range tasks {
        ts = append(ts, t)
    }
}*/
