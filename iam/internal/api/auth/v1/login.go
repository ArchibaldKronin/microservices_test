package v1

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	auth_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *apiAuth) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	login := req.GetLogin()
	pw := req.GetPassword()

	sessionID, err := a.authService.Login(ctx, login, pw)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			return nil, status.Errorf(codes.Unauthenticated, "error logging")
		case errors.Is(err, model.ErrAutheticationData):
			return nil, status.Errorf(codes.Unauthenticated, "error logging")
		default:
			return nil, status.Errorf(codes.Internal, "internal error")
		}
	}

	return &auth_v1.LoginResponse{SessionUuid: sessionID}, nil
}
