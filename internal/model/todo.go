package model

import (
	"github.com/pkg/errors"
	"time"
)

type TodoID string

var NilTodoID string

type Todo struct {
	ID TodoID `json:"id,omitempty" db:"todo_id"`
	UserID *UserID `json:"userID" db:"user_id"`
	Title *string `json:"title" db:"title"`
	Color *string `json:"color" db:"color"`
	Description *string `json:"description" db:"description"`
	IsFinished *bool `json:"isFinished" db:"is_finished"`

	CreatedAt *time.Time `json:"createdAt,omitempty" db:"created_at"`
	LastEditedAt *time.Time `json:"lastEditedAt,omitempty" db:"last_edited_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

func (a *Todo) VerifyCreate() error {
	if a.UserID == nil || len(*a.UserID) == 0 {
		return errors.New("userID is required")
	}

	if a.Title == nil || len(*a.Title) == 0 {
		return errors.New("title is required")
	}

	if a.Color == nil {
		color := "#FFFFFF"
		a.Color = &color
	}

	if a.Description == nil {
		description := ""
		a.Description = &description
	}

	return nil
}


func (a *Todo) VerifyUpdate() error {
	if a.UserID == nil || len(*a.UserID) == 0 {
		return errors.New("userID is required")
	}

	if a.Title == nil || len(*a.Title) == 0 {
		return errors.New("title is required")
	}

	if a.Color == nil || len(*a.Color) == 0 {
		return errors.New("color is required")
	}

	if a.Description == nil || len(*a.Description) == 0 {
		return errors.New("description is required")
	}

	if a.IsFinished == nil{
		return errors.New("isFinished is required")
	}

	if a.LastEditedAt == nil{
		return errors.New("lastEditedAt is required")
	}

	return nil
}

