package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/narravabrion/go-cms-server/internal/models"
)

type UserStore struct {
	db *sql.DB
}


func (us *UserStore) Create(ctx context.Context, tx *sql.Tx, user *models.User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at `
	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password.Hash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}

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
func (us *UserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {

	query := `SELECT id, username, email, password, created_at FROM users WHERE email=$1 AND is_active=true`

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	var user models.User
	err := us.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
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

func (us *UserStore) Delete(ctx context.Context, userID int64) error {
	return withTx(us.db, ctx, func(tx *sql.Tx) error {
		if err := us.delete(ctx, tx, userID); err != nil {
			return err
		}
		if err := us.deleteUserInvitation(ctx, tx, userID); err != nil {
			return err
		}
		return nil
	})
}

// check delete method conflict
func (us *UserStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// result, err := us.db.ExecContext(ctx, query, id)
	// if err != nil {
	// 	return err
	// }
	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	// 	return err
	// }
	// if rowsAffected == 0 {
	// 	return ErrNotFound
	// }

	return nil
}

func (us *UserStore) Update(ctx context.Context, user *models.User) error {
	return nil
}

func (us *UserStore) CreateAndInvite(ctx context.Context, user *models.User, token string, invitationExp time.Duration) error {
	return withTx(us.db, ctx, func(tx *sql.Tx) error {
		if err := us.Create(ctx, tx, user); err != nil {
			return err
		}
		if err := us.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (us *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExp time.Duration, userID int64) error {
	query := ` INSERT INTO user_invitations (token, user_id, expiry ) VALUES ($1, $2, $3)`
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}
	return nil
}

func (us *UserStore) Activate(ctx context.Context, token string) error {
	return withTx(us.db, ctx, func(tx *sql.Tx) error {

		user, err := us.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		user.IsActive = true
		if err := us.update(ctx, tx, user); err != nil {
			return err
		}

		if err := us.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (us *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (user *models.User, err error) {
	query := ` SELECT u.id, u.username, u.email, u.created_at, u.is_active FROM users u
 JOIN user_invitations ui on u.id = ui.user_id
 WHERE ui.token = $1 AND ui.expiry > $2`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])

	user = &models.User{}

	err = tx.QueryRowContext(ctx, query, hashedToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		switch err {

		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (us *UserStore) update(ctx context.Context, tx *sql.Tx, user *models.User) error {
	query := ` UPDATE users SET username=$1, email=$2, is_active=$3 WHERE id=$4`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

