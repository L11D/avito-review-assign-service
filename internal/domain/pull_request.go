package domain

import "time"

type PRStatus string

const (
	StatusOpen   PRStatus = "OPEN"
	StatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	Status    PRStatus   `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
	MergedAt  *time.Time `db:"merged_at"`
	AuthorID  string     `db:"author_id"`
}
