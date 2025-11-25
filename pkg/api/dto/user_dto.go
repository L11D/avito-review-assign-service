package dto

type UserSetIsActiveDTO struct {
	UserID   string `binding:"required" json:"user_id"`
	IsActive *bool  `binding:"required" json:"is_active"`
}

type UserDTO struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserPRsDTO struct {
	UserID       string                `json:"user_id"`
	PullRequests []PullRequestShortDTO `json:"pull_requests"`
}
