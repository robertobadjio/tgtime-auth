package pg_db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/robertobadjio/platform-common/pkg/db"
	"github.com/robertobadjio/tgtime-auth/internal/repository/user/model"
)

// GetUser ???
func (r *PgUserRepository) GetUser(ctx context.Context, l string) (*model.User, error) {
	builderSelect := sq.Select(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(id, login, password).
		Where(sq.Eq{"login": l})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("get user query: %w", err)
	}

	q := db.Query{
		Name:     "user_repository.GetUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return nil, nil
}
