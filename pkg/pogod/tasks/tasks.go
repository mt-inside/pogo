package tasks

import (
	"github.com/mt-inside/pogo/pkg/pogod/model"
	. "github.com/mt-inside/pogo/pkg/pogod/task"
)

/* TODO: find tasks etc in here. Easier to return gRPC errors from here,
* doesn't rely on any FSM data
 */

func Add(task *Task) {
	model.Add(task)
}

func List() []*Task {
	ts := make([]*Task, 0)
	for _, t := range model.List() {
		ts = append(ts, t)
	}
	return ts
}

func Status() (PogodState, *Task, uint32) {
	ret := post(StatusEvent, nil)
	stat := ret.(*StatusResult)
	return stat.state, stat.task, stat.remain
}

func Start(id int64) error {
	/* this is 12 kinds of nasty.
	 * Should really have a "startRet" type, which is unpacked here. That's what ret.(error) is degenerately doing
	 * But, err, as a field, would be typed, so could just be used, nil or not. Here we have to change type, and that panics if it's nil (even though we know that's ok), so we use the two-return version just to get... the zero type... which happens to be nil.
	 * Then we have to "use" ok.*/
	/* Should also be building a "StartArgs" here, but it's unary atm so meh */
	ret := post(StartEvent, id)
	err, ok := ret.(error)
	ok = ok
	return err
}

func Stop() error {
	ret := post(StopEvent, nil)
	err, ok := ret.(error)
	ok = ok
	return err
}

func Complete(id int64) error {
	ret := post(CompleteEvent, id)
	err, ok := ret.(error)
	ok = ok
	return err
}
