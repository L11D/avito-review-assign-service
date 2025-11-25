package repo

import (
	"context"

	"github.com/L11D/avito-review-assign-service/internal/domain"
	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type teamRepo struct {
	db     *sqlx.DB
	qb     sq.StatementBuilderType
	getter *trmsqlx.CtxGetter
}

func NewTeamRepo(db *sqlx.DB, getter *trmsqlx.CtxGetter) *teamRepo {
	return &teamRepo{
		db:     db,
		qb:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		getter: getter,
	}
}

func (r *teamRepo) Save(ctx context.Context, team domain.Team) (domain.Team, error) {
	query := r.qb.
		Insert("teams").
		Columns("name").
		Values(team.Name).
		Suffix("RETURNING id, name")

	sql, args, err := query.ToSql()
	if err != nil {
		return domain.Team{}, err
	}

	var createdTeam domain.Team

	err = r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &createdTeam, sql, args...)
	if err != nil {
		return domain.Team{}, err
	}

	return createdTeam, nil
}

func (r *teamRepo) GetByName(ctx context.Context, name string) (domain.Team, error) {
	query := r.qb.
		Select("id", "name").
		From("teams").
		Where(sq.Eq{"name": name})

	sql, args, err := query.ToSql()
	if err != nil {
		return domain.Team{}, err
	}

	var team domain.Team

	err = r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &team, sql, args...)
	if err != nil {
		return domain.Team{}, err
	}

	return team, nil
}

func (r *teamRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Team, error) {
	query := r.qb.
		Select("id", "name").
		From("teams").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return domain.Team{}, err
	}

	var team domain.Team

	err = r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &team, sql, args...)
	if err != nil {
		return domain.Team{}, err
	}

	return team, nil
}
