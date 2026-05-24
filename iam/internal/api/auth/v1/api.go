package v1

import (
	"github.com/ArchibaldKronin/microservices_test/iam/internal/service"
	auth_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/auth/v1"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

type apiAuth struct {
	auth_v1.UnimplementedAuthServiceServer

	authv3.UnimplementedAuthorizationServer

	authService service.AuthService
}

func NewApiAuth(authService service.AuthService) *apiAuth {
	return &apiAuth{authService: authService}
}
