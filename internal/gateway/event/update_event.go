package events

import (
	"net/http"
	"strconv"

	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

type UpdateEventRequest struct {
	Title       string   `json:"title" valid:"length(3|100)"`
	Description string   `json:"description" valid:"length(3|1000)"`
	Location    string   `json:"location" valid:"length(3|100)"`
	Category    int      `json:"category_id" valid:"range(1|8)"`
	Capacity    int      `json:"capacity" valid:"range(0|20000)"`
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
// @Param json body NewEventRequest true "Данные для обновления события"
// @Param image formData file false "Изображение события"
// @Success 200 {object} NewEventResponse "Успешное обновление события"
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

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidID)
		return
	}

	url, err := h.uploadImage(r.Context(), media, w)
	if err != nil {
		return
	}

	event := toPBEvent(req, session.UserID)
	event.ID = int32(id)
	event.Image = url

	event, err = h.EventService.UpdateEvent(r.Context(), event)
	if err != nil {
		h.deleteImage(r.Context(), url)
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.NotFound:
				utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrEventNotFound)
				return
			case grpcCodes.PermissionDenied:
				utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrAccessDenied)
				return
			case grpcCodes.Internal:
				utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
				return
			}
		}

		h.logger.Error(r.Context(), "update event", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	eventResp := eventToEventResponse(event)
	resp := NewEventResponse{
		Event: eventResp,
	}
	utils.WriteResponse(w, http.StatusOK, resp)
}
