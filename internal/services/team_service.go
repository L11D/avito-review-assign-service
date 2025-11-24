package services

import (
	"github.com/L11D/avito-review-assign-service/internal/domain"
)

type TeamRepo interface {
	Create(team domain.Team) (domain.Team, error)
}
