package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/narravabrion/go-cms-server/internal/models"
)



type RoleStore struct {
	db *sql.DB
}

func (rs *RoleStore) GetByName(ctx context.Context, roleName string) (*models.Role, error) {
	query := `SELECT id, name, level, description FROM roles WHERE name = $1`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	role := &models.Role{}
	err := rs.db.QueryRowContext(ctx, query, roleName).Scan(&role.ID, &role.Name, &role.Level, &role.Description)
	if err != nil {
		return nil, err
	}
	return role, nil
}
