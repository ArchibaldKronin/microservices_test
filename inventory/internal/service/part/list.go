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

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	if filter == nil {
		return nil, model.ErrInvalidArgument
	}
	parts, err := s.repo.ListParts(ctx, filter)
	if err != nil {

		var errConv *repoModel.MetadataParseValueError
		if errors.As(err, &errConv) && parts != nil {
			logger.Warn(
				ctx,
				"error converting Value type",
				zap.Any("metadata_value", errConv.Value),
				zap.Error(err),
			)
			return parts, nil
		}

		logger.Warn(ctx, "error listing parts", zap.Error(err))
		return nil, fmt.Errorf("error list parts service :%w", model.ErrUnexpected)
	}

	return parts, nil
}
