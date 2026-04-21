package part

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
)

func (s *service) GetPart(ctx context.Context, id string) (*model.Part, error) {
	part, err := s.repo.GetPart(ctx, id)
	if err != nil {
		if errors.Is(err, repoModel.ErrNotFound) {
			return nil, model.ErrNotFound
		}

		var errConv *repoModel.MetadataParseValueError
		if errors.As(err, &errConv) && part != nil {
			slog.Warn("error converting Value type", "metadata_value", errConv.Value, "error", err)
			return part, nil
		}

		return nil, fmt.Errorf("error get parts service :%w", model.ErrUnexpected)
	}

	return part, nil
}
