package events

import (
	"net/http"

	"kudago/internal/gateway/utils"
	httpErrors "kudago/internal/http/errors"
	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
)

// AddEvent создает новое событие в системе.
// @Summary Создание события
// @Description Создает новое событие в системе. Необходимо передать JSON-объект с данными события.
// @Tags events
// @Accept  json
// @Produce  json
// @Param json body NewEventRequest true "Данные для создания события"
// @Success 201 {object} NewEventResponse "Событие успешно создано"
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

	req, media, reqErr := parseEventData(r)
	if reqErr != nil {
		utils.WriteResponse(w, http.StatusBadRequest, reqErr)
		return
	}

	_, err := govalidator.ValidateStruct(&req)
	if err != nil {
		utils.ProcessValidationErrors(w, err)
		return
	}

	reqErr = checkNewEventRequest(req)
	if reqErr != nil {
		utils.WriteResponse(w, http.StatusBadRequest, reqErr)
		return
	}

	url, err := h.uploadImage(r.Context(), media, w)
	if err != nil {
		return
	}

	event := toPBEvent(req, session.UserID)
	event.Image = url

	event, err = h.EventService.AddEvent(r.Context(), event)
	if err != nil {
		h.deleteImage(r.Context(), url)

		switch err {
		case models.ErrInvalidCategory:
			utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidCategory)

		default:
			h.logger.Error(r.Context(), "create event", err)
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		}
		return
	}

	eventResp := eventToEventResponse(event)
	resp := NewEventResponse{
		Event: eventResp,
	}

	utils.WriteResponse(w, http.StatusOK, resp)
}
