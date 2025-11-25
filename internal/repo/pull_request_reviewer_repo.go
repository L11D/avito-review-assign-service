package repo

import (
	"context"

	"github.com/L11D/avito-review-assign-service/internal/domain"
	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
)

type pullRequestReviewerRepo struct {
	db     *sqlx.DB
	qb     sq.StatementBuilderType
	getter *trmsqlx.CtxGetter
}

func NewPullRequestReviewerRepo(db *sqlx.DB, getter *trmsqlx.CtxGetter) *pullRequestReviewerRepo {
	return &pullRequestReviewerRepo{
		db:     db,
		qb:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		getter: getter,
	}
}

func (r *pullRequestReviewerRepo) Save(
	ctx context.Context,
	prReviewer domain.PullRequestReviewer,
) (domain.PullRequestReviewer, error) {
	query := r.qb.
		Insert("pull_request_reviewers").
		Columns("pull_request_id", "user_id").
		Values(prReviewer.PullRequestID, prReviewer.UserID).
		Suffix("RETURNING pull_request_id, user_id")

	sql, args, err := query.ToSql()
	if err != nil {
		return domain.PullRequestReviewer{}, err
	}

	var createdPRReviewer domain.PullRequestReviewer

	err = r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &createdPRReviewer, sql, args...)
	if err != nil {
		return domain.PullRequestReviewer{}, err
	}

	return createdPRReviewer, nil
}

func (r *pullRequestReviewerRepo) GetPRUsersIds(ctx context.Context, prId string) ([]string, error) {
	query := r.qb.
		Select("user_id").
		From("pull_request_reviewers").
		Where(sq.Eq{"pull_request_id": prId})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var usersIds []string

	err = r.getter.DefaultTrOrDB(ctx, r.db).SelectContext(ctx, &usersIds, sql, args...)
	if err != nil {
		return nil, err
	}

	return usersIds, nil
}

func (r *pullRequestReviewerRepo) DeleteByPRAndUserId(ctx context.Context, prId string, userId string) error {
	query := r.qb.
		Delete("pull_request_reviewers").
		Where(sq.Eq{"pull_request_id": prId, "user_id": userId})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
