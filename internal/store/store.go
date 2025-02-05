package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/narravabrion/go-cms-server/internal/models"
)

var (
	ErrNotFound = errors.New("resource not found!")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *models.Post) error
		GetByID(context.Context, int64) (*models.Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *models.Post) error

	}
	Users interface {
		Create(context.Context, *models.User) error
	}
}

func NewStrorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db},
		Users: &UserStore{db},
	}
}
