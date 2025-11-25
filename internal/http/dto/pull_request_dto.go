package dto

import (
	"time"

	"github.com/L11D/avito-review-assign-service/internal/domain"
)

type PullRequestCreateDTO struct {
	ID       string `binding:"required,min=1,max=50" json:"pull_request_id"`
	Name     string `binding:"required"              json:"pull_request_name"`
	AuthorID string `binding:"required"              json:"author_id"`
}

type PullRequestDTO struct {
	ID        string          `json:"pull_request_id"`
	Name      string          `json:"pull_request_name"`
	AuthorID  string          `json:"author_id"`
	Status    domain.PRStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	MergedAt  *time.Time      `json:"merged_at,omitempty"`
	Reviewers []string        `json:"assigned_reviewers"`
}

type PullRequestShortDTO struct {
	ID       string          `json:"pull_request_id"`
	Name     string          `json:"pull_request_name"`
	AuthorID string          `json:"author_id"`
	Status   domain.PRStatus `json:"status"`
}

type PullRequestMergeDTO struct {
	ID string `binding:"required,min=1,max=50" json:"pull_request_id"`
}

type PullRequestReassignDTO struct {
	PullRequestID string `binding:"required,min=1,max=50" json:"pull_request_id"`
	OldReviewerID string `binding:"required,min=1,max=50" json:"old_reviewer_id"`
}
