package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"taskbot/cmd/web/storage"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Storage{db: db}, nil

}

func (s *Storage) Save(t *storage.Task) error {
	stmt := "INSERT INTO tasks (content, created) VALUES (?, UTC_TIMESTAMP())"

	_, err := s.db.Exec(stmt, t.Content)

	if err != nil {
		return err
	}

	return nil
}
func (s *Storage) SaveUser(u *storage.User) error {
	stmt := "INSERT INTO users (chatid, username) VALUES (?, ?)"
	_, err := s.db.Exec(stmt, u.Chatid, u.Username)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RetrieveUser(u *storage.User) (*storage.User, error) {
	stmt := "SELECT FROM users where username = ?"
	row := s.db.QueryRow(stmt, u.Username)
	var retrievedUser storage.User
	err := row.Scan(&retrievedUser.Username, &retrievedUser.Chatid)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows were returned, indicating no matching user
			return nil, nil
		}
		return nil, err
	}

	return &retrievedUser, nil
}

func (s *Storage) Tasks() ([]*storage.Task, error) {
	stmt := `SELECT id, content, created FROM tasks`

	rows, err := s.db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []*storage.Task

	for rows.Next() {
		var t *storage.Task

		err = rows.Scan(&t.ID, &t.Content, &t.Created)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, storage.ErrNoSavedTasks
	}

	return tasks, nil
}

//func (s *Storage) PastTasks() ( []*storage.Task, error ) {
//	stmt := `SELECT content, deadline, created FROM tasks
//	WHERE deadline < UTC_TIMESTAMP()`
//
//	rows, err := s.db.Query(stmt)
//	if err != nil {
//		return nil, err
//	}
//
//	defer rows.Close()
//
//	var tasks []*storage.Task
//
//	for rows.Next() {
//		var t task
//
//		err = rows.Scan(&t.Header, &t.Content, &t,Deadline, &t.Created)
//		if err != nil {
//			return nil, err
//		}
//
//		tasks = append(tasks, task)
//	}
//
//	if err = rows.Err(); err != nil {
//		return nil, err
//	}
//
//	if len(tasks) == 0 {
//		return nil, storage.ErrNoPastTasks
//	}
//
//	return tasks, nil
//}

func (s *Storage) Remove(id int) error {
	stmt := "DELETE FROM tasks WHERE id = ?"
	result, err := s.db.Exec(stmt, id)
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

func (s *Storage) Complete(id int) error {
	stmt := "UPDATE tasks SET completed = 1 WHERE id = ?"
	result, err := s.db.Exec(stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("Task wasn't completed")
	}

	return nil
}

//func (s *Storage) Clear() error {
//	stmt := "DELETE FROM tasks WHERE deadline < UTC_TIMESTAMP()"
//
//	result, err := s.db.Exec(stmt)
//
//	if err != nil {
//		return err
//	}
//
//	rowsAffected, err := result.RowsAffected()
//	if err != nil {
//		return err
//	}
//	if rowsAffected == 0 {
//		return fmt.Errorf("Tasks wasn't removed")
//	}
//
//	return nil
//}

func (s *Storage) IsExists(content string) (bool, error) {
	stmt := "SELECT * FROM tasks where content = ?"
	result, err := s.db.Exec(stmt, content)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 1 {
		return true, nil
	} else if rowsAffected == 0 {
		return false, nil
	} else {
		return false, fmt.Errorf("More than one tasks with the same data.")
	}
}
func (s *Storage) Init() error {
	stmt := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INT AUTO_INCREMENT PRIMARY KEY,
        content VARCHAR(255) UNIQUE NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        completed BOOLEAN
    );
`
	stmt2 := `
CREATE TABLE IF NOT EXISTS users (
	username VARCHAR(255) PRIMARY KEY NOT NULL,
	chatid INTEGER NOT NULL
);
`
	_, err := s.db.Exec(stmt)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(stmt2)
	if err != nil {
		return err
	}

	return nil
}
