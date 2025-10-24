package category

import (
	"context"
	"encoding/json"
	"fmt"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
)

func (r *Repository) FetchCategoryExtended(ctx context.Context, categoryID int64) (model.CategoryExtended, error) {
	q := `
        WITH recursive category_tree AS (
            SELECT id, name, parent_id
            FROM category
            WHERE id = $1
            
            UNION ALL
            
            SELECT c.id, c.name, c.parent_id
            FROM category c
            JOIN category_tree ct ON c.parent_id = ct.id
        )
        SELECT 
            c.id,
            c.name,
            p.id as "parent.id",
            p.name as "parent.name",
            COALESCE(
                (SELECT json_agg(
                    json_build_object(
                        'id', child.id,
                        'name', child.name
                    )
                ) FROM category child WHERE child.parent_id = c.id),
                '[]'::json
            ) as children
        FROM category c
        LEFT JOIN category p ON c.parent_id = p.id
        WHERE c.id = $1;
    `

	var dbResult struct {
		ID         int64           `db:"id"`
		Name       string          `db:"name"`
		ParentID   *int64          `db:"parent.id"`
		ParentName *string         `db:"parent.name"`
		Children   json.RawMessage `db:"children"`
	}

	err := r.db.GetContext(ctx, &dbResult, q, categoryID)
	if err != nil {
		return model.CategoryExtended{}, fmt.Errorf("repository.FetchCategoryExtended: %w", err)
	}

	var parent *model.Category
	if dbResult.ParentID != nil {
		parent = &model.Category{
			ID:   *dbResult.ParentID,
			Name: *dbResult.ParentName,
		}
	}

	var children []*model.Category
	if len(dbResult.Children) > 0 {
		if err = json.Unmarshal(dbResult.Children, &children); err != nil {
			return model.CategoryExtended{}, fmt.Errorf("repository.FetchCategoryExtended: failed to unmarshal children: %w", err)
		}
	}

	return model.CategoryExtended{
		ID:       dbResult.ID,
		Name:     dbResult.Name,
		Parent:   parent,
		Children: children,
	}, nil
}
