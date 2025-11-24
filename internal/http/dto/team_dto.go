package dto

type TeamDTO struct {
	Name    string          `json:"team_name" binding:"required,min=1,max=50"`
	Members []TeamMemberDTO `json:"members" binding:"required,dive"`
}

type TeamMemberDTO struct {
	Id       string `json:"user_id" binding:"required,min=1,max=50"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}
