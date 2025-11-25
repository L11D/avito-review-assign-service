package services

import (
	"context"
	"errors"
	"sort"

	"github.com/L11D/avito-review-assign-service/internal/domain"
	appErrors "github.com/L11D/avito-review-assign-service/internal/errors"
	"github.com/L11D/avito-review-assign-service/internal/http/dto"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
)

type PullRequestRepo interface {
	Save(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error)
}

type PullRequestReviewerRepo interface {
	Save(ctx context.Context, prReviewer domain.PullRequestReviewer) (domain.PullRequestReviewer, error)
}

type UserRepoPRService interface {
	GetByID(ctx context.Context, userId string) (domain.User, error)
	GetByTeamID(ctx context.Context, teamId uuid.UUID) ([]domain.User, error)
}

type pullRequestService struct {
	PRRepo      PullRequestRepo
	PRReviewerRepo PullRequestReviewerRepo
	userRepo    UserRepoPRService
	trManager *manager.Manager
}

func NewPullRequestService(prRepo PullRequestRepo, prReviewerRepo PullRequestReviewerRepo, userRepo UserRepoPRService, trManager *manager.Manager) *pullRequestService {
	return &pullRequestService{
		PRRepo:      prRepo,
		PRReviewerRepo: prReviewerRepo,
		userRepo:    userRepo,
		trManager: trManager,
	}
}

func (s *pullRequestService) Create(ctx context.Context, pr dto.PullRequestCreateDTO) (dto.PullRequestDTO, error){
	domainPR := domain.PullRequest{
		Id:       pr.Id,
		Name:     pr.Name,
		AuthorId: pr.AuthorId,
	}

	var createdPR domain.PullRequest
	var createReviewerIds []string

	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		userCreator, err := s.userRepo.GetByID(ctx, pr.AuthorId)
		if err != nil {
			if errors.Is(appErrors.MapPgError(err), appErrors.ErrNotFound) {
				return appErrors.NewNotFoundError("User with ID '" + pr.AuthorId + "'")
			}
			return err
		}

		usersInTeam, err := s.userRepo.GetByTeamID(ctx, userCreator.TeamId)
		if err != nil {
			return err
		}

		PR, err := s.PRRepo.Save(ctx, domainPR)
		if err != nil {
			if errors.Is(appErrors.MapPgError(err), appErrors.ErrAlreadyExists) {
				return appErrors.NewPullRequestExistsError(pr.Id)
			}
			return err
		}

		reviewersIds := chooseReviewers(usersInTeam, pr.AuthorId)
		for _, reviewerId := range reviewersIds {
			prReviewer := domain.PullRequestReviewer{
				PullRequestId: pr.Id,
				UserId:        reviewerId,
			}
			createdPRReviewer, err := s.PRReviewerRepo.Save(ctx, prReviewer)
			if err != nil {
				return err
			}
			createReviewerIds = append(createReviewerIds, createdPRReviewer.UserId)
		}

		createdPR = PR

		return nil
	})

	if err != nil {
		return dto.PullRequestDTO{}, err
	}

	createdPRDTO := dto.PullRequestDTO{
		Id:        createdPR.Id,
		Name:      createdPR.Name,
		AuthorId:  createdPR.AuthorId,
		Status:    createdPR.Status,
		CreatedAt: createdPR.CreatedAt,
		MergedAt:  createdPR.MergedAt,
		Reviewers: createReviewerIds,
	}

	return createdPRDTO, nil
}

func chooseReviewers(users []domain.User, creatorId string) []string {
	var probableReviewers []domain.User
	for _, user := range users {
		if user.Id != creatorId && user.IsActive {
			probableReviewers = append(probableReviewers, user)
		}
	}

	sort.Slice(probableReviewers, func(i, j int) bool {
		return probableReviewers[i].AssignRate < probableReviewers[j].AssignRate
	})

	if len(probableReviewers) > 2 {
		probableReviewers = probableReviewers[:2]
	}

	var reviewers []string
	for _, reviewer := range probableReviewers {
		reviewers = append(reviewers, reviewer.Id)
	}

	return reviewers
}