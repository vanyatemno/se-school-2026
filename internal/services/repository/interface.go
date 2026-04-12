package repository

type RepositoriesService interface {
	UpdateRepository(id uint) error
	UpdateAllRepositories() error
}
