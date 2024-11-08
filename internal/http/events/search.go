package events

import (
	"net/http"
	"strconv"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"
)

// @Summary Поиск событий
// @Description Поиск событий по ключевым словам, датам, тегам и категории
// @Tags events
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество событий на странице (по умолчанию 30)"
// @Param query query string false "Ключевые слова для поиска"
// @Param event_start query string false "Дата начала события в формате YYYY-MM-DD"
// @Param event_end query string false "Дата окончания события в формате YYYY-MM-DD"
// @Param tags query []string false "Список тегов"
// @Param category_id query int false "ID категории"
// @Success 200 {object} GetEventsResponse "Список событий"
// @Failure 400 {object} httpErrors.HttpError "Invalid Data"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events [get]
func (h EventHandler) SearchEvents(w http.ResponseWriter, r *http.Request) {
	paginationParams := utils.GetPaginationParams(r)

	query := r.URL.Query().Get("query")
	eventStart := r.URL.Query().Get("event_start")
	eventEnd := r.URL.Query().Get("event_end")
	categoryIDStr := r.URL.Query().Get("category_id")
	tags := r.URL.Query()["tags"]

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		categoryID = 0
	}

	params := models.SearchParams{
		Query:      query,
		EventStart: eventStart,
		EventEnd:   eventEnd,
		Tags:       tags,
		Category:   categoryID,
	}

	events, err := h.service.SearchEvents(r.Context(), params, paginationParams)
	if err != nil {
		h.logger.Error(r.Context(), "search", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	resp := writeEventsResponse(events, paginationParams.Limit)
	utils.WriteResponse(w, http.StatusOK, resp)
}
