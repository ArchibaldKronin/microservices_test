package part

import (
	"context"
	"slices"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	"github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/converter"
)

func (r *repository) ListParts(_ context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var parts []*model.Part

outer:
	for _, part := range r.data {

		if len(filter.Uuids) != 0 {
			if !slices.Contains(filter.Uuids, part.Uuid) {
				continue outer
			}
		}

		if len(filter.Names) != 0 {
			if !slices.Contains(filter.Names, part.Name) {
				continue outer
			}
		}

		if len(filter.Categorys) != 0 {
			categoryMatched := false
			for _, cat := range filter.Categorys {
				if part.Category == converter.CategoryToRepo(cat) {
					categoryMatched = true
					break
				}
			}
			if !categoryMatched {
				continue outer
			}
		}

		if len(filter.Countrys) != 0 {
			if !slices.Contains(filter.Countrys, part.Manufacturer.Country) {
				continue outer
			}
		}

		if len(filter.Tags) != 0 {
			tagMatched := false

		tagOuter:
			for _, reqTag := range filter.Tags {
				for _, innerTag := range part.Tags {
					if innerTag == reqTag {
						tagMatched = true
						break tagOuter
					}
				}
			}
			if !tagMatched {
				continue outer
			}
		}
		parts = append(parts, converter.PartToDomain(&part))
	}

	return parts, nil
}
