package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID         string    `db:"id"`
	Username   string    `db:"username"`
	IsActive   bool      `db:"is_active"`
	TeamID     uuid.UUID `db:"team_id"`
	AssignRate int       `db:"assign_rate"`
}
