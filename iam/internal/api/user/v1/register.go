package v1

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	user_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *apiUser) Register(ctx context.Context, req *user_v1.RegisterRequest) (*user_v1.RegisterResponse, error) {
	registrInfo := req.GetInfo()
	pw := registrInfo.GetPassword()
	userInfo := registrInfo.GetInfo()

	userID, err := a.userService.RegisterUser(ctx, model.User{
		Login: userInfo.GetLogin(),
		Email: userInfo.GetEmail(),
	}, pw)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "внутрення ошибка")
	}

	return &user_v1.RegisterResponse{
		UserUuid: userID,
	}, nil
}
