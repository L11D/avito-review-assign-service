package repo

import (
	"github.com/jmoiron/sqlx"
	sq "github.com/Masterminds/squirrel"
	"github.com/L11D/avito-review-assign-service/internal/services"
	"github.com/L11D/avito-review-assign-service/internal/domain"
)

type teamRepo struct {
	db *sqlx.DB
    qb sq.StatementBuilderType
}

func NewTeamRepo(db *sqlx.DB) services.TeamRepo {
	return &teamRepo{
		db: db,
        qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *teamRepo) Create(team domain.Team) (domain.Team, error) {
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
	err = r.db.Get(&createdTeam, sql, args...)
	if err != nil {
		return domain.Team{}, err
	}

	return createdTeam, nil
}