package domain

import (
	"github.com/L11D/avito-review-assign-service/pkg/api/dto"
	"time"
)

type PullRequest struct {
	ID        string       `db:"id"`
	Name      string       `db:"name"`
	Status    dto.PRStatus `db:"status"`
	CreatedAt time.Time    `db:"created_at"`
	MergedAt  *time.Time   `db:"merged_at"`
	AuthorID  string       `db:"author_id"`
}
