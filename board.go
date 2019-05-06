package nonota

import (
	"fmt"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type TaskTime struct {
	Start time.Time
	End   time.Time
	Note  string

	// Duration is *not* end-start; rather, it's how much work was recorded
	// within that timeframe.
	Duration time.Duration
}

type Task struct {
	Title       string
	Description string
	Times       []*TaskTime
}

func (t *Task) AddTaskTime(tt *TaskTime) {
	t.Times = append(t.Times, tt)
}

func (t *Task) TotalTime(fromTime, toTime time.Time) time.Duration {
	var total time.Duration
	for _, tt := range t.Times {
		if tt.Start.After(fromTime) && tt.End.Before(toTime) {
			total += tt.Duration
		}
	}
	return total
}

type List struct {
	Title string
	Tasks []*Task
}

func (l *List) TotalTime(fromTime, toTime time.Time) time.Duration {
	var total time.Duration
	for _, t := range l.Tasks {
		total += t.TotalTime(fromTime, toTime)
	}
	return total
}

type Board struct {
	Lists []*List
}

func (b *Board) MoveUp(list *List) {
	for i := 1; i < len(b.Lists); i++ {
		if b.Lists[i] == list {
			b.Lists[i-1], b.Lists[i] = b.Lists[i], b.Lists[i-1]
			break
		}
	}
}

func (b *Board) MoveDown(list *List) {
	for i := len(b.Lists) - 2; i >= 0; i-- {
		if b.Lists[i] == list {
			b.Lists[i+1], b.Lists[i] = b.Lists[i], b.Lists[i+1]
			break
		}
	}
}

func (b *Board) MoveTaskUp(task *Task) {
	for i := 0; i < len(b.Lists); i++ {
		for j := 0; j < len(b.Lists[i].Tasks); j++ {
			if b.Lists[i].Tasks[j] != task {
				continue
			}
			if i == 0 && j == 0 {
				// Nothing to do.
				return
			}
			if j == 0 {
				// Remove from this list and put on previous list
				b.Lists[i].Tasks = b.Lists[i].Tasks[1:]
				b.Lists[i-1].Tasks = append(b.Lists[i-1].Tasks, task)
				return
			}
			// Move to previous location within list
			b.Lists[i].Tasks[j-1], b.Lists[i].Tasks[j] = b.Lists[i].Tasks[j], b.Lists[i].Tasks[j-1]
			return
		}
	}
}

func (b *Board) MoveTaskDown(task *Task) {
	lenLists := len(b.Lists)
	for i := lenLists - 1; i >= 0; i-- {
		lenTasks := len(b.Lists[i].Tasks)
		for j := lenTasks - 1; j >= 0; j-- {
			if b.Lists[i].Tasks[j] != task {
				continue
			}
			if i == lenLists-1 && j == lenTasks-1 {
				// Nothing to do.
				return
			}
			if j == lenTasks-1 {
				// Remove from this list and put on nextlist
				var newl []*Task
				newl = append(newl, b.Lists[i].Tasks[j])
				newl = append(newl, b.Lists[i+1].Tasks...)
				b.Lists[i+1].Tasks = newl
				b.Lists[i].Tasks = b.Lists[i].Tasks[:lenTasks-1]
				return
			}
			// Move to next location within list
			b.Lists[i].Tasks[j+1], b.Lists[i].Tasks[j] = b.Lists[i].Tasks[j], b.Lists[i].Tasks[j+1]
			return
		}
	}
}

func (b *Board) AppendNewTask(list *List) {
	newTask := &Task{
		Title: "New Task",
	}
	list.Tasks = append(list.Tasks, newTask)
}

func (b *Board) AppendNewList() {
	newList := &List{
		Title: "New List",
	}
	b.Lists = append(b.Lists, newList)
}

func (b *Board) TotalTime(fromTime, toTime time.Time) time.Duration {
	var total time.Duration
	for _, l := range b.Lists {
		total += l.TotalTime(fromTime, toTime)
	}
	return total
}

func BoardFromFile(filename string) (*Board, error) {
	f, err := os.Open(filename)
	if os.IsNotExist(err) {
		return &Board{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	board := &Board{}
	dec := yaml.NewDecoder(f)
	err = dec.Decode(board)
	if err != nil {
		return nil, fmt.Errorf("error decoding file %s: %v", filename, err)
	}

	return board, nil
}

func BoardToFile(filename string, b *Board) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	return enc.Encode(b)
}
