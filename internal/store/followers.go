package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerId int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (fs *FollowerStore) Follow(ctx context.Context, FollowedID int64, userID int64) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`
	_, err := fs.db.ExecContext(ctx, query, userID, FollowedID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return ErrAlreadyFollowing
		}
		return err
	}
	return nil
}

func (fs *FollowerStore) UnFollow(ctx context.Context, FollowedID int64, userID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`
	res, err := fs.db.ExecContext(ctx, query, userID, FollowedID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFollowing
	}
	return nil
}
