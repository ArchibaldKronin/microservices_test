package part

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	if filter == nil {
		return nil, model.ErrInvalidArgument
	}
	parts, err := s.repo.ListParts(ctx, filter)
	if err != nil {

		var errConv *repoModel.MetadataParseValueError
		if errors.As(err, &errConv) && parts != nil {
			slog.Warn("error converting Value type", "metadata_value", errConv.Value, "error", err)
			return parts, nil
		}

		slog.Warn("error listing parts", "error_data", err)
		return nil, fmt.Errorf("error list parts service :%w", model.ErrUnexpected)
	}

	return parts, nil
}
