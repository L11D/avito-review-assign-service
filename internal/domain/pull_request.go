package domain

import "time"

type PullRequest struct {
	Id        string     `db:"id"`
	Name      string     `db:"name"`
	Status    string     `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
	MergedAt  *time.Time `db:"merged_at"`
	AuthorId  string     `db:"author_id"`
}
