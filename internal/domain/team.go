package domain

import (
	"github.com/google/uuid"
)

type Team struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
