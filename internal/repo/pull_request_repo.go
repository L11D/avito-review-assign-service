package repo
import (
	"context"

	"github.com/L11D/avito-review-assign-service/internal/domain"
	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
)

type pullRequestRepo struct {
	db *sqlx.DB
	qb sq.StatementBuilderType
	getter *trmsqlx.CtxGetter
}

func NewPullRequestRepo(db *sqlx.DB, getter *trmsqlx.CtxGetter) *pullRequestRepo {
	return &pullRequestRepo{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		getter: getter,
	}
}

func (r *pullRequestRepo) Save(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error) {
	query := r.qb.
		Insert("pull_requests").
		Columns("id", "name", "author_id").
		Values(pr.Id, pr.Name, pr.AuthorId).
		Suffix("RETURNING id, name, status, created_at, merged_at, author_id")

	sql, args, err := query.ToSql()
	if err != nil {
		return domain.PullRequest{}, err
	}

	var createdPR domain.PullRequest
	err = r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &createdPR, sql, args...)
	if err != nil {
		return domain.PullRequest{}, err
	}

	return createdPR, nil
}

func (r *pullRequestRepo) GetByUserId(ctx context.Context, userId string) ([]domain.PullRequest, error) {
	query := r.qb.
		Select("id", "name", "status", "created_at", "merged_at", "author_id").
		From("pull_requests").
		Join("pull_request_reviewers prr ON pull_requests.id = prr.pull_request_id").
		Where(sq.Eq{"prr.user_id": userId})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var prs []domain.PullRequest
	err = r.getter.DefaultTrOrDB(ctx, r.db).SelectContext(ctx, &prs, sql, args...)
	if err != nil {
		return nil, err
	}

	return prs, nil
}