package nonota

import (
	"fmt"
	"time"
)

type Work struct {
	Task *Task

	workTime  TaskTime
	stopChan  chan struct{}
	paused    bool
	workEnded bool
}

func NewWork(task *Task) *Work {
	return &Work{
		Task: task,
		workTime: TaskTime{
			Start: time.Now(),
		},
		stopChan: make(chan struct{}),
	}
}

func StartWork(task *Task) *Work {
	w := &Work{
		Task: task,
		workTime: TaskTime{
			Start:    time.Now(),
			Duration: time.Minute,
		},
		stopChan: make(chan struct{}),
	}

	go w.handler()

	return w
}

func (w *Work) handler() {
	ticker := time.NewTicker(time.Second)
	lastTickerStart := time.Now()

out:
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			if w.paused {
				lastTickerStart = now
				continue
			}

			w.workTime.Duration += now.Sub(lastTickerStart)
			lastTickerStart = now
		case <-w.stopChan:
			break out
		}
	}

	ticker.Stop()
}

func (w *Work) StopWork() error {
	if w.workEnded {
		return fmt.Errorf("work already stopped")
	}
	w.workEnded = true

	if w.stopChan != nil {
		close(w.stopChan)
		w.workTime.End = time.Now()
	} else {
		w.workTime.End = w.workTime.Start
	}
	w.Task.AddTaskTime(&w.workTime)
	return nil
}

func (w *Work) AdjustWorkDuration(newDuration time.Duration) {
	w.workTime.Duration = newDuration
}

func (w *Work) PauseWork() {
	w.paused = true
}

func (w *Work) ResumeWork() {
	w.paused = false
}

func (w *Work) CurrentDuration() time.Duration {
	return w.workTime.Duration
}

func (w *Work) SetNote(note string) {
	w.workTime.Note = note
}
