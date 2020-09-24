package database

import (
	"context"
	"github.com/ejiro-edwin/todolist/internal/model"
	"github.com/pkg/errors"
)

type TodoDB interface {
	CreateTodo(ctx context.Context, todo *model.Todo) error
	UpdateTodo(ctx context.Context, todo *model.Todo) error
	GetTodoByID(ctx context.Context, todoID model.TodoID) (*model.Todo, error)
	ListTodosByUserID(ctx context.Context, userID model.UserID) ([]*model.Todo, error)
	DeleteTodo(ctx context.Context, todoID model.TodoID) (bool, error)
}

const createTodoQuery = `
	INSERT INTO todo (user_id, title, color, description)
		VALUES (:user_id, :title, :color, :description)
	RETURNING todo_id;
`
func (d *database) CreateTodo(ctx context.Context, todo *model.Todo) error {
	rows, err := d.conn.NamedQueryContext(ctx, createTodoQuery, todo)
	if err != nil {
		return err
	}

	//we need return id
	defer rows.Close()
	rows.Next()
	if err := rows.Scan(&todo.ID); err != nil {
		return err
	}

	return nil
}

const updateTodoQuery = `
	UPDATE todo 
	SET title = :title,
		color = :color,
		description = :description,
		is_finished = :is_finished,
		last_edited_at = NOW()
	WHERE todo_id = :todo_id;
`
func (d *database) UpdateTodo(ctx context.Context, todo *model.Todo) error {
	result, err := d.conn.NamedExecContext(ctx, updateTodoQuery, todo)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("Todo not found")
	}

	return nil
}

const getTodoByIDQuery = `
	SELECT todo_id, user_id, title, color, description, is_finished, created_at, last_edited_at, deleted_at
	FROM todo
	WHERE todo_id = $1;
`
func (d *database) GetTodoByID(ctx context.Context, todoID model.TodoID) (*model.Todo, error) {
	var todo model.Todo
	if err := d.conn.GetContext(ctx, &todo, getTodoByIDQuery, todoID); err != nil {
		return nil, errors.Wrap(err, "could not get todo")
	}

	return &todo, nil
}

const listTodosByUserIDQuery = `
	SELECT todo_id, user_id, title, color, description, is_finished, created_at, last_edited_at, deleted_at
	FROM todo
	WHERE user_id = $1 AND deleted_at IS NULL
	ORDER BY created_at DESC;
`
func (d *database) ListTodosByUserID(ctx context.Context, userID model.UserID) ([]*model.Todo, error) {
	var todos []*model.Todo
	if err := d.conn.SelectContext(ctx, &todos, listTodosByUserIDQuery, userID); err != nil {
		return nil, errors.Wrap(err, "could not get user's todos")
	}

	return todos, nil
}

const deleteTodoQuery = `
	UPDATE todo
	SET deleted_at = NOW()
	WHERE todo_id = $1 AND deleted_at IS NULL;
`
func (d *database) DeleteTodo(ctx context.Context, todoID model.TodoID) (bool, error) {
	result, err := d.conn.ExecContext(ctx, deleteTodoQuery, todoID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}
