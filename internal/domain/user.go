package domain

import (
	"github.com/google/uuid"
)

type User struct {
	Id       string    `db:"id"`
	Username string    `db:"username" validate:"required,min=1,max=50"`
	IsActive bool      `db:"is_active"`
	TeamId   uuid.UUID `db:"team_id"`
}
