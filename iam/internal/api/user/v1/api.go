package v1

import (
	"github.com/ArchibaldKronin/microservices_test/iam/internal/service"
	user_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/user/v1"
)

type apiUser struct {
	userService service.UserService

	user_v1.UnimplementedUserServiceServer
}

func NewApiUser(userService service.UserService) *apiUser {
	return &apiUser{userService: userService}
}
