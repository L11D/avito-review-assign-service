package services

import (
	"context"
	"errors"

	"github.com/L11D/avito-review-assign-service/internal/domain"
	appErrors "github.com/L11D/avito-review-assign-service/internal/errors"
	"github.com/L11D/avito-review-assign-service/pkg/api/dto"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
)

type UserRepo interface {
	Save(ctx context.Context, user domain.User) (domain.User, error)
	GetByTeamID(ctx context.Context, teamId uuid.UUID) ([]domain.User, error)
	Update(ctx context.Context, user domain.User) (domain.User, error)
	GetByID(ctx context.Context, userId string) (domain.User, error)
}

type PullRequestRepoUserService interface {
	GetByUserId(ctx context.Context, userId string) ([]domain.PullRequest, error)
}

type TeamRepoUserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.Team, error)
}

type userService struct {
	userRepo  UserRepo
	teamRepo  TeamRepoUserService
	prRepo    PullRequestRepoUserService
	trManager *manager.Manager
}

func NewUserService(
	userRepo UserRepo,
	teamRepo TeamRepoUserService,
	prRepo PullRequestRepoUserService,
	trManager *manager.Manager,
) *userService {
	return &userService{
		userRepo:  userRepo,
		teamRepo:  teamRepo,
		prRepo:    prRepo,
		trManager: trManager,
	}
}

func (s *userService) CreateUsersInTeam(
	ctx context.Context,
	teamId uuid.UUID,
	members []dto.TeamMemberDTO,
) ([]dto.TeamMemberDTO, error) {
	createdMembers := make([]dto.TeamMemberDTO, len(members))
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		for i, member := range members {
			user := memberDTOtoUser(member, teamId)

			createdUser, err := s.userRepo.Save(ctx, user)
			if err != nil {
				if errors.Is(appErrors.MapPgError(err), appErrors.ErrAlreadyExists) {
					return appErrors.NewUserExistsError(member.ID)
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
	user, err := s.userRepo.GetByID(ctx, userSetIsActiveDTO.UserID)
	if err != nil {
		if errors.Is(appErrors.MapPgError(err), appErrors.ErrNotFound) {
			return dto.UserDTO{}, appErrors.NewNotFoundError("User with ID '" + userSetIsActiveDTO.UserID + "'")
		}

		return dto.UserDTO{}, err
	}

	if user.IsActive != *userSetIsActiveDTO.IsActive {
		user.IsActive = *userSetIsActiveDTO.IsActive

		user, err = s.userRepo.Update(ctx, user)
		if err != nil {
			return dto.UserDTO{}, err
		}
	}

	team, err := s.teamRepo.GetByID(ctx, user.TeamID)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.UserDTO{
		ID:       user.ID,
		Username: user.Username,
		IsActive: user.IsActive,
		TeamName: team.Name,
	}, err
}

func (s *userService) GetReviews(ctx context.Context, userId string) (dto.UserPRsDTO, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		if errors.Is(appErrors.MapPgError(err), appErrors.ErrNotFound) {
			return dto.UserPRsDTO{}, appErrors.NewNotFoundError("User with ID '" + userId + "'")
		}

		return dto.UserPRsDTO{}, err
	}

	prs, err := s.prRepo.GetByUserId(ctx, userId)
	if err != nil {
		return dto.UserPRsDTO{}, err
	}

	prDTOs := make([]dto.PullRequestShortDTO, len(prs))
	for i, pr := range prs {
		prDTOs[i] = dto.PullRequestShortDTO{
			ID:       pr.ID,
			Name:     pr.Name,
			Status:   pr.Status,
			AuthorID: pr.AuthorID,
		}
	}

	return dto.UserPRsDTO{
		UserID:       user.ID,
		PullRequests: prDTOs,
	}, nil
}

func (s *userService) IncrementAssignRate(ctx context.Context, userId string) (domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		if errors.Is(appErrors.MapPgError(err), appErrors.ErrNotFound) {
			return domain.User{}, appErrors.NewNotFoundError("User with ID '" + userId + "'")
		}

		return domain.User{}, err
	}

	user.AssignRate += 1

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return updatedUser, nil
}

func memberDTOtoUser(dto dto.TeamMemberDTO, teamId uuid.UUID) domain.User {
	return domain.User{
		ID:       dto.ID,
		Username: dto.Username,
		IsActive: *dto.IsActive,
		TeamID:   teamId,
	}
}

func userToMemberDTO(user domain.User) dto.TeamMemberDTO {
	return dto.TeamMemberDTO{
		ID:       user.ID,
		Username: user.Username,
		IsActive: &user.IsActive,
	}
}
