package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
)

func (s *service) GetPart(ctx context.Context, id string) (*model.Part, error) {
	part, err := s.repo.GetPart(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repoModel.ErrNotFound):
			return nil, model.ErrNotFound
		default:
			return nil, fmt.Errorf("error get parts service :%w", err)
		}
	}

	return part, nil
}
