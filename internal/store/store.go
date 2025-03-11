package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/narravabrion/go-cms-server/internal/models"
)

var (
	ErrNotFound         = errors.New("resource not found!")
	ErrAlreadyFollowing = errors.New("you are already following this user")
	ErrNotFollowing     = errors.New("you are not following this user")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *models.Post) error
		GetByID(context.Context, int64) (*models.Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *models.Post) error
		GetUserFeed(context.Context, int64, PaginationFeedQuery) ([]models.Post, error)
	}
	Users interface {
		Create(context.Context, *models.User) error
		GetByID(context.Context, int64) (*models.User, error)
		Delete(context.Context, int64) error
		Update(context.Context, *models.User) error
	}
	Followers interface {
		Follow(context.Context, int64, int64) error
		UnFollow(context.Context, int64, int64) error
	}
}

func NewStrorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Followers: &FollowerStore{db},
	}
}
