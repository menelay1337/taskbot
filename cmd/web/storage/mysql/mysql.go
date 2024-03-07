package mysql

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"taskbot/cmd/web/storage"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil

}

func (s *Storage) Save(t *storage.Task) error {
	stmt := "INSERT INTO tasks (header, content, deadline) VALUES (?, ?, ?)"
	
	result, err := s.db.Exec(stmt, t.Header, t.Content, t.Deadline)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Tasks() ( []*storage.Task, error ) {
	stmt := `SELECT header, content, deadline, created FROM tasks
	WHERE deadline > UTC_TIMESTAMP()`

	rows, err := s.db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []*storage.Task

	for rows.Next() {
		var t task

		err = rows.Scan(&t.Header, &t.Content, &t,Deadline, &t.Created)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}


func (s *Storage) PastTasks() ( []*storage.Task, error ) {
	stmt := `SELECT header, content, deadline, created FROM tasks
	WHERE deadline < UTC_TIMESTAMP()`

	rows, err := s.db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []*storage.Task

	for rows.Next() {
		var t task

		err = rows.Scan(&t.Header, &t.Content, &t,Deadline, &t.Created)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Storage) Remove(t *storage.Task) error {
	stmt := "DELETE FROM tasks WHERE header = ? AND content = ? AND deadline = ?"
	result, err := s.db.Exec(stmt, t.Header, t.Content, t.Deadline)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("Task wasn't removed")
	}

	return nil
}

func (s *Storage) Clear() error {
	stmt := "DELETE FROM tasks WHERE deadline < UTC_TIMESTAMP()"

	result, err := s.db.Exec(stmt)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Tasks wasn't removed")
	}

	return nil
}

func (s *Storage) isExists(t *storage.Task) (bool, error) {
	stmt :=	"SELECT * FROM tasks where header = ? AND content = ? AND deadline = ?"
	result, err := s.db.Exec(stmt, t.Header, t.Content, t.Deadline)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 1 {
		return true, nil
	} else uf rowsAffected == 0 {
		return false, nil
	} else {
		return false, fmt.Errorf("More than one tasks with the same data.")
	}
}


