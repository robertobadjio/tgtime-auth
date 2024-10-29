package pg_db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/robertobadjio/platform-common/pkg/db"

	"github.com/robertobadjio/tgtime-auth/internal/repository/user/model"
)

// GetUser ???
func (r *PgUserRepository) GetUser(ctx context.Context, email string) (*model.User, error) {
	builderSelect := sq.
		Select(idColumnName, emailColumnName, passwordColumnName, roleColumnName).
		PlaceholderFormat(sq.Dollar).
		From(userTableName).
		Where(sq.Eq{emailColumnName: email}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build user query: %w", err)
	}

	q := db.Query{
		Name:     "user_repository.GetUser",
		QueryRaw: query,
	}

	var user model.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	return &user, nil
}
