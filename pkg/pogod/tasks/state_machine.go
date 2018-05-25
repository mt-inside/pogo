package tasks

import (
	"log"
	"time"

	"github.com/spf13/viper"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mt-inside/pogo/pkg/pogod/model"
	. "github.com/mt-inside/pogo/pkg/pogod/task"
)

type PogodState int /* FSM: S */

const (
	Idle PogodState = iota
	RunningTask
	RunningBreak
)

type ClientEvent int /* FSM: Σ */

const (
	StartEvent ClientEvent = iota
	StopEvent
	CompleteEvent
	TimerEvent
	StatusEvent
)

type PogodData struct {
	Task      *Task
	Timer     *time.Timer
	StartTime time.Time
	EndTime   time.Time
}

func (d PogodData) addTimer() PogodData {
	if (d.Timer != nil || d.StartTime != time.Time{} || d.EndTime != time.Time{}) {
		panic("Should not be a timer")
	}

	pomoDuration := time.Duration(viper.GetInt("pomodoro_time")) * time.Minute
	d.Timer = time.NewTimer(pomoDuration)
	d.StartTime = time.Now()
	d.EndTime = d.StartTime.Add(pomoDuration)
	return d
}
func (d PogodData) clearTimer() PogodData {
	if (d.Timer == nil || d.StartTime == time.Time{} || d.EndTime == time.Time{}) {
		panic("Should be a timer")
	}

	d.Timer = nil
	d.StartTime = time.Time{}
	d.EndTime = time.Time{}
	return d
}
func (d PogodData) addTask(t *Task) PogodData {
	if d.Task != nil {
		panic("Should not be a task")
	}

	d.Task = t
	return d
}
func (d PogodData) changeTask(t *Task) PogodData {
	if d.Task == nil {
		panic("Should be a task")
	}

	d.Task = t
	return d
}
func (d PogodData) clearTask() PogodData {
	d.Task = nil
	return d
}

type taskMachine struct {
	state  PogodState
	data   *PogodData
	inbox  chan actorMessage
	output chan interface{}
}

var (
	/* FSM: s_0 */
	fsm taskMachine = taskMachine{Idle, &PogodData{}, make(chan actorMessage), make(chan interface{})}
)

/* These get state so that they can use that args to make it obvious when they want to stay in the same state, and so handler functions can be re-used */
/* I think the returning of random interfaces is the best we can do since we can't emit events (though there could totally be a generic interface that had a method that accepted n "events") */
/* It's not possible to mark struct members const, so:
* - data is passed by value so that mutating the fields in place makes no difference and doesn't accidentally change the state of the system
* - as an orthogonal problem in immutability, we'd like to force people to make new objects rather than mutating the ones they're given. This is possible by moving the type to another package, making the fields private, and adding getters.
*   - TODO: this ^^, one day */
type action func(PogodState, PogodData, interface{}) (PogodState, *PogodData, interface{})

/* Avoiding an init loop */
var (
	actions [3][5]action /* FSM: δ : S × Σ → S */
)

func init() {
	/* TODO could allow nil entries for unhandles states (print that you're doing nothing), but it adds code complexity */
	actions = [3][5]action{
		{idle_start, idle_stop, idle_complete, idle_timer, any_status},
		{task_start, task_stop, task_complete, task_timer, any_status},
		{break_start, break_stop, break_complete, break_timer, any_status},
	}

	go actorLoop()
}

/* TODO: make a method of fsm */
type actorMessage struct {
	event ClientEvent
	args  interface{}
}

func post(event ClientEvent, args interface{}) (result interface{}) {
	fsm.inbox <- actorMessage{event, args}
	result = <-fsm.output
	return
}
func actorLoop() {
	var result interface{}
	for {
		msg := <-fsm.inbox

		log.Printf("ACTION! State %d; Event %d", fsm.state, msg.event)
		fsm.state, fsm.data, result = actions[fsm.state][msg.event](fsm.state, *fsm.data, msg.args)

		fsm.output <- result
	}
}

type StatusResult struct {
	state  PogodState
	task   *Task
	remain uint32
}

func any_status(state PogodState, data PogodData, args interface{}) (PogodState, *PogodData, interface{}) {
	ret := &StatusResult{
		state,
		data.Task,
		uint32(data.EndTime.Sub(time.Now()).Seconds()),
	}

	return state, &data, ret
}

func idle_start(state PogodState, data PogodData, args interface{}) (PogodState, *PogodData, interface{}) {
	id := args.(int64)
	t, found := model.Find(id)
	if !found {
		return Idle, &data, status.Errorf(codes.NotFound, "No such task")
	}
	if t.State == DONE {
		return Idle, &data, status.Errorf(codes.FailedPrecondition, "Task already complete")
	}

	data = data.addTask(t)

	data = data.addTimer()

	go func() {
		<-data.Timer.C
		/* Stop whatever's running at the time.
		 * Recall: closes over this */
		log.Printf("Timer expired")
		post(TimerEvent, nil)
	}()

	log.Printf("Started task %v", t)

	return RunningTask, &data, nil
}
func idle_stop(state PogodState, data PogodData, args interface{}) (PogodState, *PogodData, interface{}) {
	log.Println("Nothing to do")

	return state, &data, nil
}
func idle_complete(state PogodState, data PogodData, args interface{}) (PogodState, *PogodData, interface{}) {
	id := args.(int64)
	t, found := model.Find(id)
	if !found {
		return Idle, &data, status.Errorf(codes.NotFound, "No such task")
	}
	if t.State == DONE {
		return Idle, &data, status.Errorf(codes.FailedPrecondition, "Task already completed")
	}
	t.State = DONE
	// data doesn't have a task to clear

	log.Printf("Completed task %v", t)

	return state, &data, nil
}
func idle_timer(state PogodState, data PogodData, args interface{}) (PogodState, *PogodData, interface{}) {
	panic("timer in idle mode")
}

func task_start(state PogodState, data PogodData, args interface{}) (PogodState, *PogodData, interface{}) {
	id := args.(int64)
	t, found := model.Find(id)
	if !found {
		return state, &data, status.Errorf(codes.NotFound, "No such task")
	}
	if t.State == DONE {
		return state, &data, status.Errorf(codes.FailedPrecondition, "Task already completed")
	}
	// tasks don't change state
	data = data.changeTask(t)

	log.Printf("Overwrote task with %v", t)

	return state, &data, nil
}
func task_stop(state PogodState, data PogodData, args interface{}) (PogodState, *PogodData, interface{}) {
	log.Printf("FYI running task is %v", data.Task)
	// task doesn't change state
	data = data.clearTask()

	/* TODO: makes this goroutine hang forever?
	* Thus will leak them.
	* Instead, have that goroutine select on a cancel channel as well?
	 */
	if !data.Timer.Stop() {
		/* This catches the rare case that we asked the timer to stop
		 * just as it was expiring. I.e. it had already marked itself
		 * expired, but we are sure that the expiration event hasn't
		 * yet been received. In this case we can be sure of that, because
		 * that would have moved the state machine to state Idle.
		 * Thus, we drain the channel here to make sure the other
		 * listener isn't spuriously woken up later. */
		log.Println("timer had already expired - rare")
		<-data.Timer.C /* Drain to avoid spurious wakeup */
	} else {
		log.Println("timer had not expired")
	}

	data = data.clearTimer()

	log.Println("Stopped task due to user")

	return Idle, &data, nil
}
func task_complete(state PogodState, data PogodData, args interface{}) (PogodState, *PogodData, interface{}) {
	id := args.(int64)
	t, found := model.Find(id)
	if !found {
		return state, &data, status.Errorf(codes.NotFound, "No such task")
	}
	if t.State == DONE {
		return state, &data, status.Errorf(codes.FailedPrecondition, "Task already completed")
	}

	state = RunningTask

	/* If it's that task that's running, stop the task */
	if data.Task == t {
		data.Task = nil

		if !data.Timer.Stop() {
			/* See comments in task_stop. */
			log.Println("timer had already expired - rare")
			<-data.Timer.C /* Drain to avoid spurious wakeup */
		} else {
			log.Println("timer had not expired")
		}

		data = data.clearTimer()

		state = Idle
	}

	/* Complete the task no matter what */
	t.State = DONE

	log.Printf("Completed task %v", t)

	return state, &data, nil
}
func task_timer(state PogodState, data PogodData, args interface{}) (PogodState, *PogodData, interface{}) {
	log.Printf("FYI running task was %v", data.Task)
	// task doesn't change state
	data = data.clearTask()

	data = data.clearTimer()

	log.Println("Stopped task due to timer")

	return Idle, &data, nil
}

func break_start(state PogodState, data PogodData, args interface{}) (PogodState, *PogodData, interface{}) {
	panic("no breaks")
}
func break_stop(state PogodState, data PogodData, args interface{}) (PogodState, *PogodData, interface{}) {
	panic("no breaks")
}
func break_complete(state PogodState, data PogodData, args interface{}) (PogodState, *PogodData, interface{}) {
	panic("no breaks")
}
func break_timer(state PogodState, data PogodData, args interface{}) (PogodState, *PogodData, interface{}) {
	panic("no breaks")
}
