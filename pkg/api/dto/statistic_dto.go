package dto

type UserStatisticDTO struct {
	UserID     string `json:"user_id"`
	AssignRate int    `json:"assign_rate"`
}

type AllUsersStatisticDTO struct {
	Statistics []UserStatisticDTO `json:"user_statistics"`
}
