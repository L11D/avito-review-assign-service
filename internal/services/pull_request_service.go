package services

import (
	"context"
	"errors"
	"slices"
	"sort"
	"time"

	"github.com/L11D/avito-review-assign-service/internal/domain"
	appErrors "github.com/L11D/avito-review-assign-service/internal/errors"
	"github.com/L11D/avito-review-assign-service/internal/http/dto"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
)

const MAX_REVIEWERS_PER_PR = 2

type PullRequestRepo interface {
	Save(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error)
	Update(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error)
	GetByID(ctx context.Context, prId string) (domain.PullRequest, error)
}

type PullRequestReviewerRepo interface {
	Save(ctx context.Context, prReviewer domain.PullRequestReviewer) (domain.PullRequestReviewer, error)
	GetPRUsersIds(ctx context.Context, prId string) ([]string, error)
	DeleteByPRAndUserId(ctx context.Context, prId string, userId string) error
}

type UserRepoPRService interface {
	GetByID(ctx context.Context, userId string) (domain.User, error)
	GetByTeamID(ctx context.Context, teamId uuid.UUID) ([]domain.User, error)
}

type UserServicePRService interface {
	IncrementAssignRate(ctx context.Context, userId string) (domain.User, error)
}

type pullRequestService struct {
	PRRepo         PullRequestRepo
	PRReviewerRepo PullRequestReviewerRepo
	userRepo       UserRepoPRService
	userService    UserServicePRService
	trManager      *manager.Manager
}

func NewPullRequestService(
	prRepo PullRequestRepo,
	prReviewerRepo PullRequestReviewerRepo,
	userRepo UserRepoPRService,
	userService UserServicePRService,
	trManager *manager.Manager,
) *pullRequestService {
	return &pullRequestService{
		PRRepo:         prRepo,
		PRReviewerRepo: prReviewerRepo,
		userRepo:       userRepo,
		userService:    userService,
		trManager:      trManager,
	}
}

func (s *pullRequestService) Create(ctx context.Context, pr dto.PullRequestCreateDTO) (dto.PullRequestDTO, error) {
	domainPR := domain.PullRequest{
		ID:       pr.ID,
		Name:     pr.Name,
		AuthorID: pr.AuthorID,
	}

	var (
		createdPR         domain.PullRequest
		createReviewerIds []string
	)

	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		PR, err := s.PRRepo.Save(ctx, domainPR)
		if err != nil {
			if errors.Is(appErrors.MapPgError(err), appErrors.ErrAlreadyExists) {
				return appErrors.NewPullRequestExistsError(pr.ID)
			}

			return err
		}

		reviewersIds, err := s.getReviewsForUserPR(ctx, pr.AuthorID, []string{})
		if err != nil {
			return err
		}

		for _, reviewerId := range reviewersIds {
			prReviewer := domain.PullRequestReviewer{
				PullRequestID: pr.ID,
				UserID:        reviewerId,
			}

			createdPRReviewer, err := s.PRReviewerRepo.Save(ctx, prReviewer)
			if err != nil {
				return err
			}

			createReviewerIds = append(createReviewerIds, createdPRReviewer.UserID)
		}

		createdPR = PR

		return nil
	})
	if err != nil {
		return dto.PullRequestDTO{}, err
	}

	return prToDTO(createdPR, createReviewerIds), nil
}

func (s *pullRequestService) Merge(ctx context.Context, prId string) (dto.PullRequestDTO, error) {
	pr, err := s.PRRepo.GetByID(ctx, prId)
	if err != nil {
		if errors.Is(appErrors.MapPgError(err), appErrors.ErrNotFound) {
			return dto.PullRequestDTO{}, appErrors.NewNotFoundError("Pull Request with ID '" + prId + "'")
		}

		return dto.PullRequestDTO{}, err
	}

	var returnedReviewerIds []string

	returnedPr := pr

	if pr.Status != domain.StatusMerged {
		returnedPr, returnedReviewerIds, err = s.doMerge(ctx, pr)
		if err != nil {
			return dto.PullRequestDTO{}, err
		}
	} else {
		returnedReviewerIds, err = s.PRReviewerRepo.GetPRUsersIds(ctx, pr.ID)
		if err != nil {
			return dto.PullRequestDTO{}, err
		}
	}

	return prToDTO(returnedPr, returnedReviewerIds), nil
}

func (s *pullRequestService) Reassign(
	ctx context.Context, 
	reassignDTO dto.PullRequestReassignDTO,
) (dto.PullRequestDTO, error) {
	pr, err := s.PRRepo.GetByID(ctx, reassignDTO.PullRequestID)
	if err != nil {
		if errors.Is(appErrors.MapPgError(err), appErrors.ErrNotFound) {
			return dto.PullRequestDTO{}, appErrors.NewNotFoundError("Pull Request with ID '" + reassignDTO.PullRequestID + "'")
		}

		return dto.PullRequestDTO{}, err
	}

	if pr.Status == domain.StatusMerged {
		return dto.PullRequestDTO{}, appErrors.NewPullRequestMergedError()
	}

	hasReviewer, err := s.prHasReviewer(ctx, pr.ID, reassignDTO.OldReviewerID)
	if err != nil {
		return dto.PullRequestDTO{}, err
	}

	if !hasReviewer {
		return dto.PullRequestDTO{}, appErrors.NewNotAssignedError()
	}

	newReviewersIds, err := s.getReviewsForUserPR(ctx, pr.AuthorID, []string{reassignDTO.OldReviewerID})
	if err != nil {
		return dto.PullRequestDTO{}, err
	}

	if len(newReviewersIds) < 1 {
		return dto.PullRequestDTO{}, appErrors.NewNoCandidateError()
	}

	err = s.doReassign(ctx, pr.ID, newReviewersIds[0], reassignDTO.OldReviewerID)
	if err != nil {
		return dto.PullRequestDTO{}, err
	}

	reviewerIds, err := s.PRReviewerRepo.GetPRUsersIds(ctx, pr.ID)
	if err != nil {
		return dto.PullRequestDTO{}, err
	}

	return prToDTO(pr, reviewerIds), nil
}

func (s *pullRequestService) doReassign(
	ctx context.Context,
	prId string,
	newReviewerId string,
	oldReviewerId string,
) error {
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		err := s.PRReviewerRepo.DeleteByPRAndUserId(ctx, prId, oldReviewerId)
		if err != nil {
			return err
		}

		prReviewer := domain.PullRequestReviewer{
			PullRequestID: prId,
			UserID:        newReviewerId,
		}

		_, err = s.PRReviewerRepo.Save(ctx, prReviewer)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (s *pullRequestService) doMerge(
	ctx context.Context,
	notMergedPr domain.PullRequest,
) (domain.PullRequest, []string, error) {
	notMergedPr.Status = domain.StatusMerged
	now := time.Now().UTC()
	notMergedPr.MergedAt = &now

	var (
		returnedReviewerIds []string
		returnedPr          domain.PullRequest
	)

	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		updatedPR, err := s.PRRepo.Update(ctx, notMergedPr)
		if err != nil {
			return err
		}

		returnedPr = updatedPR

		reviewersIds, err := s.PRReviewerRepo.GetPRUsersIds(ctx, returnedPr.ID)
		if err != nil {
			return err
		}

		for _, reviewerId := range reviewersIds {
			_, err = s.userService.IncrementAssignRate(ctx, reviewerId)
			if err != nil {
				return err
			}
		}

		returnedReviewerIds = reviewersIds

		return nil
	})
	if err != nil {
		return returnedPr, nil, err
	}

	return returnedPr, returnedReviewerIds, nil
}

func (s *pullRequestService) getReviewsForUserPR(
	ctx context.Context,
	userId string,
	excludeIds []string,
) ([]string, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		if errors.Is(appErrors.MapPgError(err), appErrors.ErrNotFound) {
			return nil, appErrors.NewNotFoundError("User with ID '" + userId + "'")
		}

		return nil, err
	}

	usersInTeam, err := s.userRepo.GetByTeamID(ctx, user.TeamID)
	if err != nil {
		return nil, err
	}

	excludeIds = append(excludeIds, userId)
	reviewersIds := chooseReviewers(usersInTeam, excludeIds)

	return reviewersIds, nil
}

func (s *pullRequestService) prHasReviewer(ctx context.Context, prId string, reviewerId string) (bool, error) {
	usersIds, err := s.PRReviewerRepo.GetPRUsersIds(ctx, prId)
	if err != nil {
		return false, err
	}

	if slices.Contains(usersIds, reviewerId) {
		return true, nil
	}

	return false, nil
}

func chooseReviewers(users []domain.User, excludeIds []string) []string {
	var probableReviewers []domain.User

	for _, user := range users {
		if user.IsActive {
			excluded := slices.Contains(excludeIds, user.ID)

			if !excluded {
				probableReviewers = append(probableReviewers, user)
			}
		}
	}

	sort.Slice(probableReviewers, func(i, j int) bool {
		return probableReviewers[i].AssignRate < probableReviewers[j].AssignRate
	})

	if len(probableReviewers) > MAX_REVIEWERS_PER_PR {
		probableReviewers = probableReviewers[:MAX_REVIEWERS_PER_PR]
	}

	var reviewers = []string{}
	for _, reviewer := range probableReviewers {
		reviewers = append(reviewers, reviewer.ID)
	}

	return reviewers
}

func prToDTO(pr domain.PullRequest, reviewerIds []string) dto.PullRequestDTO {
	return dto.PullRequestDTO{
		ID:        pr.ID,
		Name:      pr.Name,
		AuthorID:  pr.AuthorID,
		Status:    pr.Status,
		CreatedAt: pr.CreatedAt,
		MergedAt:  pr.MergedAt,
		Reviewers: reviewerIds,
	}
}
