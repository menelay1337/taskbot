package storage

import (
	"time"
	"errors"
)

type Storage interface {
	Save(t *Task) error
	Tasks() ([]*Task, error)
	//PastTasks() ([]*Task, error)
	//SetDeadline(days int) error
	Complete(id int) error
	Remove(id int) error
	//Clear() error
	IsExists(content string) (bool, error)
}

var ErrNoSavedTasks = errors.New("No saved tasks")
//var ErrNoPastTasks = errors.New("No saved tasks")

type Task struct {
	ID		  int
	Content   string
	//Deadline time.Time
	Created	  time.Time
	Completed bool
}


