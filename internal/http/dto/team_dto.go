package dto

type TeamDTO struct {
	Name    string          `binding:"required,min=1,max=50" json:"team_name"`
	Members []TeamMemberDTO `binding:"required,dive"         json:"members"`
}

type TeamMemberDTO struct {
	ID       string `binding:"required,min=1,max=50" json:"user_id"`
	Username string `binding:"required"              json:"username"`
	IsActive bool   `binding:"required"              json:"is_active"`
}
