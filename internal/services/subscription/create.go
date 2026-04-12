package subscription

import (
	"context"
	"errors"
	"fmt"
	"se-school/internal/models"
	"se-school/internal/models/dto"
	"se-school/internal/notifications/templates"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *Service) Create(
	ctx context.Context,
	req *dto.CreateSubscriptionRequest,
) error {
	err := s.createNewSubscription(ctx, req)
	if err != nil {
		zap.L().Error("failed to init new subscription", zap.Error(err))
	}
	err = s.createAndSendConfirmationCode(req.Email)
	if err != nil {
		zap.L().Error("failed to init new subscription", zap.Error(err))
	}

	return nil
}

func (s *Service) createNewSubscription(ctx context.Context, req *dto.CreateSubscriptionRequest) error {
	parsedRepoValues, err := parseRepoFields(req.Repo)
	if err != nil {
		zap.L().Error("failed to parse repo fields", zap.Error(err))
		return err
	}

	repo, err := s.getOrCreateRepository(ctx, parsedRepoValues)
	if err != nil {
		zap.L().Error("failed to get or create repository", zap.Error(err))
	}

	unsubCode, err := s.codesRepository.Create(models.CodeTypeUnsubscribe)
	if err != nil {
		zap.L().Error("failed to create unsubscribe code", zap.Error(err))
		return err
	}

	sub := &models.Subscription{
		RepositoryID:      repo.ID,
		UnsubscribeCodeID: unsubCode.ID,
		Email:             req.Email,
		LastSeenTag:       repo.Version,
	}
	err = s.subscriptionsRepository.Create(sub)
	if err != nil {
		zap.L().Error("failed to create subscription", zap.Error(err))
		return err
	}

	return nil
}

func parseRepoFields(repo string) (*parsedRepoValue, error) {
	values := strings.Split(repo, "/")
	if len(values) != 2 {
		return nil, fmt.Errorf("failed to parse repo fields: %s", repo)
	}

	return &parsedRepoValue{
		Owner:          values[0],
		RepositoryName: values[1],
	}, nil
}

func (s *Service) getOrCreateRepository(ctx context.Context, values *parsedRepoValue) (*models.Repository, error) {
	repository, err := s.repositoriesRepository.Find(&models.Repository{
		Owner: values.Owner,
		Name:  values.RepositoryName,
	})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		repository, err = s.createRepository(ctx, values)
		if err != nil {
			return nil, err
		}
	}

	return repository, nil
}

func (s *Service) createRepository(ctx context.Context, values *parsedRepoValue) (*models.Repository, error) {
	currentVersion, err := s.githubIntegration.GetRepositoryVersion(ctx, values.Owner, values.RepositoryName)
	if err != nil {
		return nil, err
	}
	repo := &models.Repository{
		Owner:   values.Owner,
		Name:    values.RepositoryName,
		Version: currentVersion,
	}
	err = s.repositoriesRepository.Create(repo)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (s *Service) createAndSendConfirmationCode(email string) error {
	subCode, err := s.codesRepository.Create(models.CodeTypeConfirm)
	if err != nil {
		return err
	}

	err = s.notificationService.SendEmail(
		[]string{email},
		templates.Confirmation,
		templates.ConfirmEmailPayload{
			Code: subCode.Code,
		})
	if err != nil {
		return err
	}

	return nil
}
