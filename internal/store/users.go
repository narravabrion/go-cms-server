package store

import (
	"context"
	"database/sql"

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
