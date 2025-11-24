package repo

import (
	"context"

	"github.com/L11D/avito-review-assign-service/internal/domain"
	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

func (r *userRepo) Save(ctx context.Context, user domain.User) (domain.User, error) {
	query := r.qb.
		Insert("users").
		Columns("id", "username", "is_active", "team_id").
		Values(user.Id, user.Username, user.IsActive, user.TeamId).
		Suffix("RETURNING id, username, team_id")

	sql, args, err := query.ToSql()
	if err != nil {
		return domain.User{}, err
	}

	var createdUser domain.User
	err = r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &createdUser, sql, args...)
	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}

func (r *userRepo) GetByTeamID(ctx context.Context, teamId uuid.UUID) ([]domain.User, error) {
	query := r.qb.
		Select("id", "username", "is_active", "team_id").
		From("users").
		Where(sq.Eq{"team_id": teamId})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var users []domain.User
	err = r.getter.DefaultTrOrDB(ctx, r.db).SelectContext(ctx, &users, sql, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}
