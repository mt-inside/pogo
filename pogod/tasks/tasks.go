package tasks

import (
	"log"
)

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

var state Pogod = Pogod{Idle, nil}

func Start(task *Task) {
	log.Println("lololol")
	//TODO
	//change state.
	//timer.
}

//TODO: state machine for the whole thing, that's what starts a task. in
//TODO: there needs to be an actor. A single source of truth and mediator
//layer, that will e.g. complete task 1 when 2 is running, but not 2.
//layer1: grpc server - unmarshals args to internal types (e.g. protoTask)
//layer2: tasks actor: mediates everything, finds existing objects, etc.
//Talks to task db, timers, audit db, main state machine (which is a db)
