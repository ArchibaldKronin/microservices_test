package part

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	"github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/converter"
	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	query := bson.M{}

	if filter != nil {
		if len(filter.Uuids) != 0 {
			query["uuid"] = bson.M{"$in": filter.Uuids}
		}
		if len(filter.Names) != 0 {
			query["name"] = bson.M{"$in": filter.Names}
		}

		if len(filter.Categorys) != 0 {
			query["category"] = bson.M{"$in": filter.Categorys}
		}

		if len(filter.Countrys) != 0 {
			query["manufacturer.country"] = bson.M{"$in": filter.Countrys}
		}
		if len(filter.Tags) != 0 {
			query["tags"] = bson.M{"$all": filter.Tags}
		}
	}

	cursor, err := r.data.Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error find part by filter %#v: %w", filter, err)
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			slog.Warn("failed to close cursor", "close_cursor_error", cerr)
		}
	}()

	var repoParts []repoModel.Part
	err = cursor.All(ctx, &repoParts)
	if err != nil {
		return nil, fmt.Errorf("error reading parts: %w", err)
	}

	parts := make([]*model.Part, 0, len(repoParts))
	var converterError *repoModel.MetadataParseValueError

	for _, rp := range repoParts {
		dp, err := converter.PartToDomain(&rp)
		if err != nil {

			var convErr *repoModel.MetadataParseValueError
			if errors.As(err, &convErr) {
				if converterError == nil {
					converterError = convErr
				}

				// slog.Warn("error converting Value type", "metadata_value", convErr.Value)
			} else {
				return nil, fmt.Errorf("error converter unexpected: %w", err)
			}
		}

		parts = append(parts, dp)
	}

	if converterError != nil {
		return parts, converterError
	}

	return parts, nil
}
