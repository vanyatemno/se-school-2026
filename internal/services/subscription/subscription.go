package subscription

import (
	"se-school/internal/integrations/github"
	"se-school/internal/notifications"
	"se-school/internal/repositories/code"
	"se-school/internal/repositories/repository"
	"se-school/internal/repositories/subscription"
)

type Service struct {
	subscriptionsRepository subscription.SubscriptionsRepository
	repositoriesRepository  repository.RepositoriesRepository
	codesRepository         code.CodesRepository
	githubIntegration       github.GithubIntegration
	notificationService     notifications.NotificationsService
}

func New(
	subscriptionsRepository subscription.SubscriptionsRepository,
	repositoriesRepository repository.RepositoriesRepository,
	codesRepository code.CodesRepository,
	githubIntegration github.GithubIntegration,
	notificationService notifications.NotificationsService,
) *Service {
	return &Service{
		subscriptionsRepository: subscriptionsRepository,
		repositoriesRepository:  repositoriesRepository,
		codesRepository:         codesRepository,
		githubIntegration:       githubIntegration,
		notificationService:     notificationService,
	}
}
