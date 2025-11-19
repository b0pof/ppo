package category

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/internal/model"
	"github.com/b0pof/ppo/internal/util/http/response"
)

func (h *Category) GetApi1CategoriesId(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	categoryID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		response.BadRequest(w, "category id parameter is required")
		return
	}

	categoryExtended, err := h.category.FetchCategoryExtended(ctx, categoryID)
	if err != nil {
		h.log.Warn("failed to get category", err)
		response.Internal(w, err)
		return
	}

	response.OK(w, toCategoryDTO(categoryExtended))
}

func toCategoryDTO(category model.CategoryExtended) dto.GetCategoryResponse {
	var parent *dto.CategoryInfo
	if category.Parent != nil {
		parent = &dto.CategoryInfo{
			Id:   category.Parent.ID,
			Name: category.Parent.Name,
		}
	}

	children := make([]dto.CategoryInfo, 0, len(category.Children))
	for _, child := range category.Children {
		if child == nil {
			continue
		}

		children = append(children, dto.CategoryInfo{
			Id:   child.ID,
			Name: child.Name,
		})
	}

	return dto.GetCategoryResponse{
		Id:       category.ID,
		Name:     category.Name,
		Parent:   parent,
		Children: children,
	}
}
