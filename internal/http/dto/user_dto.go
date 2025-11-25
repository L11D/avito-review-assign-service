package dto

type UserSetIsActiveDTO struct {
	UserId   string `json:"user_id" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

type UserDTO struct {
	Id       string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserPRsDTO struct {
	UserId       string                `json:"user_id"`
	PullRequests []PullRequestShortDTO `json:"pull_requests"`
}
