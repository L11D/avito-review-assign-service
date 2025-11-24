package services

import (
	"context"

	"github.com/L11D/avito-review-assign-service/internal/domain"
	"github.com/L11D/avito-review-assign-service/internal/http/dto"
	"github.com/google/uuid"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type TeamRepo interface {
	Save(ctx context.Context, team domain.Team) (domain.Team, error)
}

type UserService interface {
	CreateUsersInTeam(ctx context.Context, teamId uuid.UUID, members []dto.TeamMemberDTO) ([]dto.TeamMemberDTO, error)
}

type teamService struct {
	repo TeamRepo
	userService UserService
	trManager *manager.Manager
}

func NewTeamService(repo TeamRepo, userService UserService, trManager *manager.Manager) *teamService {
	return &teamService{
		repo: repo,
		userService: userService,
		trManager: trManager,
	}
}

func (s *teamService) Create(ctx context.Context, team dto.TeamDTO) (dto.TeamDTO, error) {
	domainTeam := domain.Team{
		Name: team.Name,
	}

	var createdTeamDTO dto.TeamDTO
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		createdTeam, err := s.repo.Save(ctx, domainTeam)
		if err != nil {
			return err
		}

		createdMembers, err := s.userService.CreateUsersInTeam(ctx, createdTeam.Id, team.Members)
		if err != nil {
			return err
		}

		createdTeamDTO = dto.TeamDTO {
			Name: createdTeam.Name,
			Members: createdMembers,
		}

		return nil
	})
	
	return createdTeamDTO, err
}
