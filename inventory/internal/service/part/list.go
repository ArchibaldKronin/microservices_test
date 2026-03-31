package part

import (
	"context"
	"fmt"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	if filter == nil {
		return nil, model.ErrInvalidArgument
	}
	parts, err := s.repo.ListParts(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error list parts service :%w", err)

	}

	return parts, nil
}
