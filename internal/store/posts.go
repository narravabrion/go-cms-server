package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/narravabrion/go-cms-server/internal/models"
)

type Post struct {
}

type PostStore struct {
	db *sql.DB
}

func (ps *PostStore) Create(ctx context.Context, post *models.Post) error {
	query := `INSERT INTO posts (title, content, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := ps.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (ps *PostStore) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	query := `SELECT id, user_id, title, content, created_at, updated_at, tags FROM posts WHERE id=$1`
	var post models.Post
	err := ps.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (ps *PostStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id=$1`
	// var post models.Post

	result, err := ps.db.ExecContext(ctx, query, id)

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

func (ps *PostStore) Update(ctx context.Context, post *models.Post) error {
	query := `UPDATE posts SET title =$1, content = $2 WHERE id= $3`
	_, err := ps.db.ExecContext(ctx, query, post.Title, post.Content, post.ID)
	if err != nil {
		return err
	}
	return nil
}
