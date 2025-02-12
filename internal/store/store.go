package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/narravabrion/go-cms-server/internal/models"
)

var (
	ErrNotFound         = errors.New("resource not found!")
	ErrAlreadyFollowing = errors.New("you are already following this user")
	ErrNotFollowing     = errors.New("you are not following this user")
	ErrDuplicateEmail  = errors.New("email already exists")
	ErrDuplicateUsername  = errors.New("username already exists")
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
		Create(context.Context, *sql.Tx, *models.User) error
		GetByID(context.Context, int64) (*models.User, error)
		Delete(context.Context, int64) error
		Update(context.Context, *models.User) error
		CreateAndInvite(context.Context, *models.User, string, time.Duration) error
		Activate(context.Context, string) error
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

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
