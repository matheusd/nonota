package nonota

type User struct {
	CurrentWorks []*Work
}

func (u *User) StopWork(work *Work) error {
	workIdx := -1
	for i, w := range u.CurrentWorks {
		if w == work {
			workIdx = i
			break
		}
	}

	if workIdx != -1 {
		u.CurrentWorks = append(u.CurrentWorks[:workIdx], u.CurrentWorks[workIdx+1:]...)
	}

	return work.StopWork()
}

func (u *User) ExcludeWork(work *Work) {
	workIdx := -1
	for i, w := range u.CurrentWorks {
		if w == work {
			workIdx = i
			break
		}
	}

	if workIdx != -1 {
		u.CurrentWorks = append(u.CurrentWorks[:workIdx], u.CurrentWorks[workIdx+1:]...)
	}
}

// ToggleWorkOnTask either resumes (or starts) working on the given task or
// pauses it (if already working).
func (u *User) ToggleWorkOnTask(task *Task) *Work {
	var work *Work
	for _, w := range u.CurrentWorks {
		if w.Task == task {
			if w.paused {
				w.ResumeWork()
			} else {
				w.PauseWork()
			}
			work = w
		} else {
			w.PauseWork()
		}
	}

	if work == nil {
		work = StartWork(task)
		u.CurrentWorks = append(u.CurrentWorks, work)
	} else if work.paused {
		return nil
	}

	return work
}

func (u *User) WorkForTask(task *Task) *Work {
	for _, w := range u.CurrentWorks {
		if w.Task == task {
			return w
		}
	}
	return nil
}
