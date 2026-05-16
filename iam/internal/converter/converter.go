package converter

import (
	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	common_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/common/v1"
)

func UserToProto(user model.User) *common_v1.User {
	return &common_v1.User{
		UserUuid:            user.UserUUID,
		Login:               user.Login,
		Email:               user.Email,
		NotificationMethods: notificationMethodsToProro(user.NotificationMethods),
	}
}

func notificationMethodsToProro(nm []model.NotificationMethod) []*common_v1.NotificationMethod {
	nmProto := make([]*common_v1.NotificationMethod, 0)
	if len(nm) == 0 {
		return nmProto
	}

	for _, method := range nm {
		protoMethod := &common_v1.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		}
		nmProto = append(nmProto, protoMethod)
	}

	return nmProto
}
