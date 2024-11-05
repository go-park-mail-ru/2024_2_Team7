package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type UpdateEventRequest struct {
	Title       string   `json:"title" valid:"length(3|50), omitempty"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	Category    int      `json:"category_id"`
	Capacity    int      `json:"capacity"`
	Tag         []string `json:"tag"`
	EventStart  string   `json:"event_start" valid:"rfc3339"`
	EventEnd    string   `json:"event_end" valid:"rfc3339"`
}

// UpdateEvent обновляет данные существующего события.
// @Summary Обновление события
// @Description Обновляет данные существующего события. Необходимо передать JSON-объект с данными события и идентификатором события в URL.
// @Tags events
// @Accept  json
// @Produce  json
// @Param id path int true "Идентификатор события"
// @Param json body UpdateEventRequest true "Данные для обновления события"
// @Param image formData file false "Изображение события"
// @Success 200 {object} EventResponse "Успешное обновление события"
// @Failure 400 {object} httpErrors.HttpError "Неверные данные"
// @Failure 401 {object} httpErrors.HttpError "Неавторизован"
// @Failure 403 {object} httpErrors.HttpError "Доступ запрещен"
// @Failure 404 {object} httpErrors.HttpError "Событие не найдено"
// @Failure 500 {object} httpErrors.HttpError "Внутренняя ошибка сервера"
// @Router /events/{id} [put]
func (h EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req UpdateEventRequest
	jsonData := r.FormValue("json")
	err = json.Unmarshal([]byte(jsonData), &req)
	if err != nil {
		fmt.Println(err)
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}

	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		utils.ProcessValidationErrors(w, err)
		return
	}

	event := models.Event{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		EventStart:  req.EventStart,
		EventEnd:    req.EventEnd,
		AuthorID:    session.UserID,
		Tag:         req.Tag,
		Location:    req.Location,
		CategoryID:  req.Category,
		Capacity:    req.Capacity,
	}

	media, err := utils.HandleImageUpload(r)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	event, err = h.service.UpdateEvent(r.Context(), event, media)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEventNotFound):
			utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrEventNotFound)
		case errors.Is(err, models.ErrAccessDenied):
			utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrAccessDenied)
		default:
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		}
		return
	}
	resp := eventToEventResponse(event)
	utils.WriteResponse(w, http.StatusOK, resp)
}
