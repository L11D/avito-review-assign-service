package domain

type PullRequestReviewer struct {
	PullRequestID string `db:"pull_request_id"`
	UserID        string `db:"user_id"`
}
