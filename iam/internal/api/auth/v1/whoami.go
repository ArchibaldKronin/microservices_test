package v1

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/converter"
	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	auth_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *apiAuth) Whoami(ctx context.Context, req *auth_v1.WhoamiRequest) (*auth_v1.WhoamiResponse, error) {
	sessionID := req.GetSessionUuid()

	user, err := a.authService.Whoami(ctx, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			return nil, status.Errorf(codes.NotFound, "error session not found")
		case errors.Is(err, model.ErrAutheticationData):
			return nil, status.Errorf(codes.Unauthenticated, "error logging")
		default:
			return nil, status.Errorf(codes.Internal, "internal error")
		}
	}

	return &auth_v1.WhoamiResponse{User: converter.UserToProto(*user)}, nil
}
