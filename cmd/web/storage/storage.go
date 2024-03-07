package storage

import (
	"time"
)

type Storage interface {
	Save(t *Task) error
	Tasks() ([]*Task, error)
	PastTasks() ([]*Task, error)
	Remove(t *Task) error
	Clear() error
	IsExists() (bool, error)
}

type Task struct {
	Header string
	Content  string
	Deadline time.Time
	Created	 time.Time
}
