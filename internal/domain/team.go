package domain

import (
    "github.com/google/uuid"
)

type Team struct {
    Id   uuid.UUID `db:"id"`
    Name string    `db:"name" validate:"required,min=1,max=50"`
}