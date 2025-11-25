package services

import (
	"context"
	"sort"

	"github.com/L11D/avito-review-assign-service/internal/domain"
	"github.com/L11D/avito-review-assign-service/pkg/api/dto"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type UserRepoStatistic interface {
	GetAll(ctx context.Context) ([]domain.User, error)
}

type statisticService struct {
	userRepo  UserRepoStatistic
	trManager *manager.Manager
}

func NewStatisticService(
	userRepo UserRepoStatistic,
	trManager *manager.Manager,
) *statisticService {
	return &statisticService{
		userRepo:  userRepo,
		trManager: trManager,
	}
}

func (s *statisticService) GetUsersStatistic(ctx context.Context) (dto.AllUsersStatisticDTO, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return dto.AllUsersStatisticDTO{}, err
	}

	userStatistics := make([]dto.UserStatisticDTO, 0, len(users))
	for _, user := range users {
		userStatistics = append(userStatistics, dto.UserStatisticDTO{
			UserID:     user.ID,
			AssignRate: user.AssignRate,
		})
	}

	sort.Slice(userStatistics, func(i, j int) bool {
		return userStatistics[i].AssignRate > userStatistics[j].AssignRate
	})

	return dto.AllUsersStatisticDTO{
		Statistics: userStatistics,
	}, nil
}
