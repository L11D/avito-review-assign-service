package domain

type PullRequestReviewer struct {
	PullRequestId string `db:"pull_request_id"`
	UserId        string `db:"user_id"`
}
