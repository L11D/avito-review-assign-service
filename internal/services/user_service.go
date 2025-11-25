package services

import (
	"context"
	"errors"

	"github.com/L11D/avito-review-assign-service/internal/domain"
	appErrors "github.com/L11D/avito-review-assign-service/internal/errors"
	"github.com/L11D/avito-review-assign-service/internal/http/dto"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
)

type UserRepo interface {
	Save(ctx context.Context, user domain.User) (domain.User, error)
	GetByTeamID(ctx context.Context, teamId uuid.UUID) ([]domain.User, error)
	SetIsActive(ctx context.Context, userId string, isActive bool) (domain.User, error)
}

type TeamRepoUserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.Team, error)
}

type userService struct {
	userRepo UserRepo
	teamRepo TeamRepoUserService
	trManager   *manager.Manager
}

func NewUserService(userRepo UserRepo, teamRepo TeamRepoUserService, trManager *manager.Manager) *userService {
	return &userService{userRepo: userRepo, teamRepo: teamRepo, trManager: trManager}
}

func (s *userService) CreateUsersInTeam(ctx context.Context, teamId uuid.UUID, members []dto.TeamMemberDTO) ([]dto.TeamMemberDTO, error) {
	createdMembers := make([]dto.TeamMemberDTO, len(members))
	
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		for i, member := range members {
			user := memberDTOtoUser(member, teamId)
			createdUser, err := s.userRepo.Save(ctx, user)
			if err != nil {
				if errors.Is(appErrors.MapPgError(err), appErrors.ErrAlreadyExists) {
					return appErrors.NewUserExistsError(member.Id)
				}
				return err
			}
			createdMembers[i] = userToMemberDTO(createdUser)
		}
		return nil
	})
	

	return createdMembers, err
}

func (s *userService) GetTeamMembers(ctx context.Context, teamId uuid.UUID) ([]dto.TeamMemberDTO, error) {
	users, err := s.userRepo.GetByTeamID(ctx, teamId)
	if err != nil {
		return nil, err
	}

	memberDTOs := make([]dto.TeamMemberDTO, len(users))
	for i, user := range users {
		memberDTOs[i] = userToMemberDTO(user)
	}

	return memberDTOs, nil
}

func (s *userService) SetIsActive(ctx context.Context, userSetIsActiveDTO dto.UserSetIsActiveDTO) (dto.UserDTO, error) {
	updatedUser, err := s.userRepo.SetIsActive(ctx, userSetIsActiveDTO.UserId, *userSetIsActiveDTO.IsActive)
	if err != nil {
		if errors.Is(appErrors.MapPgError(err), appErrors.ErrNotFound) {
			return dto.UserDTO{}, appErrors.NewNotFoundError("User with ID '" + userSetIsActiveDTO.UserId)
		}
	}
	team, err := s.teamRepo.GetByID(ctx, updatedUser.TeamId)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.UserDTO{
		Id:       updatedUser.Id,
		Username: updatedUser.Username,
		IsActive: updatedUser.IsActive,
		TeamName: team.Name,
	}, err
}

func memberDTOtoUser(dto dto.TeamMemberDTO, teamId uuid.UUID) domain.User {
	return domain.User{
		Id: dto.Id,
		Username: dto.Username,
		IsActive: dto.IsActive,
		TeamId: teamId,
	}
}

func userToMemberDTO(user domain.User) dto.TeamMemberDTO {
	return dto.TeamMemberDTO{
		Id: user.Id,
		Username: user.Username,
		IsActive: user.IsActive,
	}
}