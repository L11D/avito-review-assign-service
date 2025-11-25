package dto

import (
	"time"

	"github.com/L11D/avito-review-assign-service/internal/domain"
)

type PullRequestCreateDTO struct {
	Id       string `json:"pull_request_id" binding:"required,min=1,max=50"`
	Name     string `json:"pull_request_name" binding:"required"`
	AuthorId string `json:"author_id" binding:"required"`
}

type PullRequestDTO struct {
	Id        string     `json:"pull_request_id"`
	Name      string     `json:"pull_request_name"`
	AuthorId  string     `json:"author_id"`
	Status    domain.PRStatus     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	MergedAt  *time.Time `json:"merged_at,omitempty"`
	Reviewers []string   `json:"assigned_reviewers"`
}

type PullRequestShortDTO struct {
	Id       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorId string `json:"author_id"`
	Status   domain.PRStatus `json:"status"`
}

type PullRequestMergeDTO struct {
	Id string `json:"pull_request_id" binding:"required,min=1,max=50"`
}

