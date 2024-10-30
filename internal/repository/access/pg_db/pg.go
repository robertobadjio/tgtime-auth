package pg_db

import (
	"github.com/robertobadjio/platform-common/pkg/db"
	"github.com/robertobadjio/tgtime-auth/internal/repository/access"
)

const accessTableName = "public.access"

const (
	endpointColumnName = "endpoint"
	roleColumnName     = "role"
)

// PgUserRepository ???
type PgUserRepository struct {
	db db.Client
}

// NewPgRepository ???
func NewPgRepository(db db.Client) access.Repository {
	return &PgUserRepository{db: db}
}
