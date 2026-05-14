package converter

import (
	"encoding/json"
	"fmt"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/iam/internal/repository/model"
)

func UserToRepo(user model.User, passwordHash ...string) (repoModel.User, error) {
	result := repoModel.User{
		UserUUID: user.UserUUID,
		Login:    user.Login,
		Email:    user.Email,
	}

	if len(passwordHash) != 0 {
		result.PasswordHash = passwordHash[0]
	}

	nm, err := json.Marshal(user.NotificationMethods)
	if err != nil {
		return result, fmt.Errorf("Error marshaling Notification Method: %w", err)
	}

	result.NotificationMethods = nm
	return result, nil
}

func UserToDomain(user repoModel.User) (model.User, error) {
	result := model.User{
		UserUUID: user.UserUUID,
		Login:    user.Login,
		Email:    user.Email,
	}

	var nm []model.NotificationMethod
	err := json.Unmarshal(user.NotificationMethods, &nm)
	if err != nil {
		return result, fmt.Errorf("Error unmarshaling Notification Method: %w", err)
	}

	result.NotificationMethods = nm
	return result, nil
}
