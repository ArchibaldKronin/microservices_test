package v1

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/converter"
	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	user_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *apiUser) GetUser(ctx context.Context, req *user_v1.GetUserRequest) (*user_v1.GetUserResponse, error) {
	userID := req.GetUserUuid()
	user, err := a.userService.GetUserByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			return nil, status.Errorf(codes.NotFound, "пользователь с id: %s не найден", userID)
		default:
			return nil, status.Errorf(codes.Internal, "internal error")
		}
	}

	return &user_v1.GetUserResponse{
		User: converter.UserToProto(*user),
	}, nil
}
