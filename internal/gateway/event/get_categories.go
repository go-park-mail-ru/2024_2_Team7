package events

import (
	"net/http"

	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"
	"kudago/internal/models"
)

// @Summary Получить все категории
// @Description Получить список всех доступных категорий событий
// @Tags categories
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Category "Список категорий"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /categories [get]
func (h EventHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.EventService.GetCategories(r.Context(), nil)
	if err != nil {
		h.logger.Error(r.Context(), "get categories", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	convertedCategories := make([]models.Category, 0, len(categories.Categories))
	for _, cat := range categories.Categories {
		convertedCategories = append(convertedCategories, models.Category{
			ID:   int(cat.ID),
			Name: cat.Name,
		})
	}

	resp := GetCategoriesResponse{
		Categories: convertedCategories,
	}
	utils.WriteResponse(w, http.StatusOK, resp)
}
