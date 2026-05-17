package user

import (
	"github.com/ArchibaldKronin/microservices_test/iam/internal/repository"
	def "github.com/ArchibaldKronin/microservices_test/iam/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) *service {
	return &service{repo: repo}
}
