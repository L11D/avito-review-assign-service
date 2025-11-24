package services

import (
	"github.com/L11D/avito-review-assign-service/internal/domain"
	"github.com/L11D/avito-review-assign-service/internal/http/dto"
)

type TeamRepo interface {
	Save(team domain.Team) (domain.Team, error)
}

type teamService struct {
	repo TeamRepo
}

func NewTeamService(repo TeamRepo) *teamService {
	return &teamService{
		repo: repo,
	}
}

func (s *teamService) Create(team dto.TeamDTO) (dto.TeamDTO, error) {
	domainTeam := domain.Team{
		Name: team.Name,
	}

	createdTeam, err := s.repo.Save(domainTeam)
	if err != nil {
		return dto.TeamDTO{}, err
	}

	return dto.TeamDTO{
		Name: createdTeam.Name,
	}, nil
}
