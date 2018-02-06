package tasks

import (
	"log"
	"time"

	"github.com/mt-inside/pogo/pogod/data"
	. "github.com/mt-inside/pogo/pogod/task"
)

/* TODO: find tasks etc in here
* factor out a separate state machine file which has an actual 2d
* activation table thing so that you can't make the wrong transitions
 */
type PogodState int

const (
	Idle PogodState = iota
	RunningTask
	RunningBreak
)

type Pogod struct {
	State PogodState
	Task  *Task
}

var pogod *Pogod = &Pogod{Idle, nil}

func State() (s PogodState, task *Task, time uint32) {
	s = pogod.State
	task = pogod.Task
	time = 0 //TODO

	return
}

func Add(title string) {
	data.Add(title)
}

func List() []*Task {
	ts := make([]*Task, 0)
	for _, t := range data.List() {
		ts = append(ts, t)
	}
	return ts
}

func Start(id int64) {
	t := data.Find(id)

	pogod.State = RunningTask
	pogod.Task = t

	timer := time.NewTimer(5 * time.Second)
	go func() {
		<-timer.C
		stop(t)
	}()

	log.Printf("Started task %v", t)
}

func stop(t *Task) {
	pogod.State = Idle
	pogod.Task = nil
	log.Printf("Stopped task %v", t)
}

func Complete(id int64) {
	t := data.Find(id)

	t.State = Done

	log.Printf("Completed task %v", t)
}

//TODO: state machine for the whole thing, that's what starts a task. in
//TODO: there needs to be an actor. A single source of truth and mediator
//layer, that will e.g. complete task 1 when 2 is running, but not 2.
//layer1: grpc server - unmarshals args to internal types (e.g. protoTask)
//layer2: tasks actor: mediates everything, finds existing objects, etc.
//Talks to task db, timers, audit db, main state machine (which is a db)
