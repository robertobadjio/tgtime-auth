package pg_db

import (
	"github.com/robertobadjio/platform-common/pkg/db"
	"github.com/robertobadjio/tgtime-auth/internal/repository/user"
)

const tableName = "user"

const (
	id       = "id"
	login    = "login"
	password = "password"
)

// PgUserRepository ???
type PgUserRepository struct {
	db db.Client
}

// NewPgRepository ???
func NewPgRepository(db db.Client) user.Repository {
	return &PgUserRepository{db: db}
}
