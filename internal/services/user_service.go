package services

import (
	"context"

	"github.com/L11D/avito-review-assign-service/internal/domain"
	"github.com/L11D/avito-review-assign-service/internal/http/dto"
	"github.com/google/uuid"
)

type UserRepo interface {
	SaveMany(ctx context.Context, users []domain.User) ([]domain.User, error)
}

type userService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *userService {
	return &userService{repo: repo}
}

func (s *userService) CreateUsersInTeam(ctx context.Context, teamId uuid.UUID, members []dto.TeamMemberDTO) ([]dto.TeamMemberDTO, error) {
	users := make([]domain.User, len(members))
	for i, member := range members {
		users[i] = memberDTOtoUser(member, teamId)
	}

	createdUsers, err := s.repo.SaveMany(ctx, users)
	if err != nil {
		return nil, err
	}

	createdMembers := make([]dto.TeamMemberDTO, len(createdUsers))
	for i, user := range createdUsers {
		createdMembers[i] = userToMemberDTO(user)
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