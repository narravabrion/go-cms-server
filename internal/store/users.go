package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/narravabrion/go-cms-server/internal/models"
)

type UserStore struct {
	db *sql.DB
}

func (us *UserStore) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at `
	err := us.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.Username,
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err

	}
	return nil
}

func (us *UserStore) GetByID(ctx context.Context, id int64) (*models.User, error) {

	query := `SELECT id, username, email FROM users WHERE id=$1`

	var user models.User
	err := us.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (us *UserStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := us.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (us *UserStore) Update(ctx context.Context, user *models.User) error {
	return nil
}
