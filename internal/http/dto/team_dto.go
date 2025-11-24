package dto

type TeamDTO struct {
	Name    string          `json:"team_name" binding:"required"`
	Members []TeamMemberDTO `json:"members" binding:"required,dive"`
}

type TeamMemberDTO struct {
	Id       string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}
