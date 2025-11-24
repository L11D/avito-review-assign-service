package services

import (
	"context"
	"errors"

	appErrors "github.com/L11D/avito-review-assign-service/internal/errors"
	"github.com/L11D/avito-review-assign-service/internal/domain"
	"github.com/L11D/avito-review-assign-service/internal/http/dto"
	"github.com/google/uuid"
)

type UserRepo interface {
	Save(ctx context.Context, user domain.User) (domain.User, error)
}

type userService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *userService {
	return &userService{repo: repo}
}

func (s *userService) CreateUsersInTeam(ctx context.Context, teamId uuid.UUID, members []dto.TeamMemberDTO) ([]dto.TeamMemberDTO, error) {
	createdMembers := make([]dto.TeamMemberDTO, len(members))
	
	for i, member := range members {
		user := memberDTOtoUser(member, teamId)
		createdUser, err := s.repo.Save(ctx, user)
		if err != nil {
			if errors.Is(appErrors.MapPgError(err), appErrors.ErrAlreadyExists) {
				return nil, appErrors.NewUserExistsError(member.Id)
			}
			return nil, err
		}
		createdMembers[i] = userToMemberDTO(createdUser)
	}

	return createdMembers, nil
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