package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/narravabrion/go-cms-server/internal/models"
)

// type Post struct {
// }

type PostStore struct {
	db *sql.DB
}

func (ps *PostStore) Create(ctx context.Context, post *models.Post) error {
	query := `INSERT INTO posts (title, content, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
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
	query := `SELECT id, user_id, title, content, created_at, updated_at, tags, version FROM posts WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var post models.Post
	err := ps.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

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
	query := `UPDATE posts SET title =$1, content = $2, version = version + 1 WHERE id= $3 AND version = $4 RETURNING version`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := ps.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(
		&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func (ps *PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginationFeedQuery) ([]models.Post, error) {
	log.Print("get user feed")
	query := `
	
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username 
		FROM posts p 
		LEFT JOIN users u on p.user_id = u.id 
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE 
			f.user_id = $1 AND 
			(p.tags @> $2 OR $2 = '{}')
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $3 OFFSET $4
		`

	rows, err := ps.db.QueryContext(ctx, query, userID, pq.Array(fq.Tags), fq.Limit, fq.Offset)
	log.Printf("rows: %+v", rows)
	log.Printf("err: %d", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []models.Post
	var user models.User

	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Version,
			pq.Array(&post.Tags),
			&user.Username,
		)
		if err != nil {
			return nil, err
		}
		feed = append(feed, post)
	}

	return feed, nil
}
