package pg_db

import (
	"github.com/robertobadjio/platform-common/pkg/db"
	"github.com/robertobadjio/tgtime-auth/internal/repository/user"
)

const userTableName = "public.user"

const (
	idColumnName       = "id"
	emailColumnName    = "email"
	passwordColumnName = "password"
	roleColumnName     = "role"
)

// PgUserRepository ???
type PgUserRepository struct {
	db db.Client
}

// NewPgRepository ???
func NewPgRepository(db db.Client) user.Repository {
	return &PgUserRepository{db: db}
}
