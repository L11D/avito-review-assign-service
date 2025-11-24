package repo

import (
	"context"

	"github.com/L11D/avito-review-assign-service/internal/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
)

type userRepo struct {
	db *sqlx.DB
	qb sq.StatementBuilderType
	getter *trmsqlx.CtxGetter
}

func NewUserRepo(db *sqlx.DB, getter *trmsqlx.CtxGetter) *userRepo {
	return &userRepo{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		getter: getter,
	}
}

func (r *userRepo) SaveMany(ctx context.Context, users []domain.User) ([]domain.User, error) {
	query := r.qb.
		Insert("users").
		Columns("id", "username", "is_active", "team_id")

	for _, user := range users {
		query = query.Values(user.Id, user.Username, user.IsActive, user.TeamId)
	}

	query = query.Suffix("RETURNING id, username, team_id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var createdUsers []domain.User
	err = r.getter.DefaultTrOrDB(ctx, r.db).SelectContext(ctx, &createdUsers, sql, args...)
	if err != nil {
		return nil, err
	}

	return createdUsers, nil
}