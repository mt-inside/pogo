package model

import (
	"sync"

	. "github.com/mt-inside/pogo/pogod/task"
)

var (
	/* Locks the store. We don't want returned object being mutated in
	 * parallel either, and we don't want multiple things interleaving
	 * mutations either. In short, we want to give out the Tasks with a
	 * lease (like a Rust borrowCount?). I.e. each task has a r/w lock and
	 * returning a copy for read Rlocks that lock until the ownership is
	 * given up */
	lock *sync.RWMutex = &sync.RWMutex{}

	tasks    map[int64]*Task = make(map[int64]*Task)
	next_idx int64           = 0
)

func Add(title string) {
	lock.Lock()
	defer lock.Unlock()

	t := NewTask(next_idx, title)
	tasks[next_idx] = t

	next_idx += 1
}

func List() map[int64]*Task {
	lock.RLock()
	defer lock.RUnlock()

	return tasks
}

func Find(id int64) (task *Task, found bool) {
	lock.RLock()
	defer lock.RUnlock()

	task, found = tasks[id]
	return
}
