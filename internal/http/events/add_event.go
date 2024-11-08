package events

import (
	"encoding/json"
	"net/http"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
)

type AddEventRequest struct {
	Title       string   `json:"title" valid:"required,length(3|100)"`
	Description string   `json:"description" valid:"required"`
	Location    string   `json:"location"`
	Category    int      `json:"category_id" valid:"required"`
	Capacity    int      `json:"capacity"`
	Tag         []string `json:"tag"`
	EventStart  string   `json:"event_start" valid:"rfc3339,required"`
	EventEnd    string   `json:"event_end" valid:"rfc3339,required"`
}

type CreateEventResponse struct {
	Event EventResponse `json:"event"`
}

// AddEvent создает новое событие в системе.
// @Summary Создание события
// @Description Создает новое событие в системе. Необходимо передать JSON-объект с данными события.
// @Tags events
// @Accept  json
// @Produce  json
// @Param json body AddEventRequest true "Данные для создания события"
// @Success 201 {object} CreateEventResponse "Событие успешно создано"
// @Failure 400 {object} httpErrors.HttpError "Неверные данные"
// @Failure 401 {object} httpErrors.HttpError "Неавторизован"
// @Failure 500 {object} httpErrors.HttpError "Внутренняя ошибка сервера"
// @Router /events [post]
func (h EventHandler) AddEvent(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	var req AddEventRequest
	jsonData := r.FormValue("json")
	err := json.Unmarshal([]byte(jsonData), &req)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}

	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		utils.ProcessValidationErrors(w, err)
		return
	}

	event := models.Event{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		EventStart:  req.EventStart,
		EventEnd:    req.EventEnd,
		AuthorID:    session.UserID,
		CategoryID:  req.Category,
		Capacity:    req.Capacity,
		Tag:         req.Tag,
	}

	media, err := utils.HandleImageUpload(r)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	event, err = h.service.AddEvent(r.Context(), event, media)
	if err != nil {
		switch err {
		case models.ErrInvalidCategory:
			utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidCategory)

		///TODO пока оставлю так, когда будет более четкая бд и ошибки для обработки, поправлю
		default:
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		}
		return
	}
	resp := eventToEventResponse(event)
	utils.WriteResponse(w, http.StatusOK, resp)
}
