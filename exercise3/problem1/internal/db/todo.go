package db

import (
	"database/sql"
	"errors"
	"time"
)

type TodoModel struct {
	ID          int
	Description string
	Done        bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m *Model) SelectTodo(id int) (*TodoModel, error) {
	statement := "SELECT id, description, done FROM todo WHERE id=$1"
	row := m.db.QueryRow(statement, id)
	todo := &TodoModel{}

	if err := row.Scan(&todo.ID, &todo.Description, &todo.Done); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return todo, nil
		}
		return nil, err
	}

	return todo, nil
}

func (m *Model) SelectTodos() ([]*TodoModel, error) {
	statement := "SELECT id, description, done FROM todo"
	rows, err := m.db.Query(statement)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	todos := make([]*TodoModel, 0)

	for rows.Next() {
		todo := &TodoModel{}
		if err := rows.Scan(&todo.ID, &todo.Description, &todo.Done); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (m *Model) InsertTodo(todo *TodoModel) error {
	statement := "INSERT INTO todo (description, done)  VALUES ($1, $2)"
	_, err := m.db.Exec(statement, todo.Description, todo.Done)
	if err != nil {
		return err
	}
	return nil
}
