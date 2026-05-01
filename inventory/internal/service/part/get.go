package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) GetPart(ctx context.Context, id string) (*model.Part, error) {
	part, err := s.repo.GetPart(ctx, id)
	if err != nil {
		if errors.Is(err, repoModel.ErrNotFound) {
			return nil, model.ErrNotFound
		}

		var errConv *repoModel.MetadataParseValueError
		if errors.As(err, &errConv) && part != nil {
			logger.Warn(
				ctx,
				"error converting Value type",
				zap.Any("metadata_value", errConv.Value),
				zap.Error(err),
			)
			return part, nil
		}

		return nil, fmt.Errorf("error get parts service :%w", model.ErrUnexpected)
	}

	return part, nil
}
