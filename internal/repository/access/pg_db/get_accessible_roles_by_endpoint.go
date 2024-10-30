package pg_db

import (
	"context"
	"fmt"

	"github.com/robertobadjio/tgtime-auth/internal/repository/access/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/robertobadjio/platform-common/pkg/db"
)

func (r *PgUserRepository) GetAccessibleRolesByEndpoint(
	ctx context.Context,
	endpoint string,
) ([]string, error) {
	builderSelect := sq.
		Select(roleColumnName).
		PlaceholderFormat(sq.Dollar).
		From(accessTableName).
		Where(sq.Eq{endpointColumnName: endpoint}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build access query: %w", err)
	}

	q := db.Query{
		Name:     "access_repository.GetAccessibleEndpointsByRole",
		QueryRaw: query,
	}

	var roles model.Roles
	err = r.db.DB().ScanOneContext(ctx, &roles, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	return roles.Roles, nil
}
