package events

import (
	"context"
	"net/http"
	"time"

	pbEvent "kudago/internal/event/api"
	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"
	pbNtf "kudago/internal/notification/api"

	"github.com/asaskevich/govalidator"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
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
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.InvalidArgument:
				utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrInvalidData)
				return
			}
		}

		h.logger.Error(r.Context(), "add event", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	eventResp := eventToEventResponse(event)
	resp := NewEventResponse{
		Event: eventResp,
	}

	err = h.sendCreatedNotifications(r.Context(), int(event.ID), session.UserID)
	if err != nil {
		h.logger.Error(r.Context(), "send create notifications", err)
	}

	utils.WriteResponse(w, http.StatusOK, resp)
}

func (h EventHandler) sendCreatedNotifications(ctx context.Context, eventID int, userID int) error {
	idsResp, err := h.EventService.GetSubscribersIDs(ctx, &pbEvent.GetSubscribersIDsRequest{UserID: int32(userID)})
	if err != nil {
		return err
	}

	req := &pbNtf.CreateNotificationsRequest{
		UserIDs: make([]int32, 0, len(idsResp.IDs)),
		Notification: &pbNtf.Notification{
			Message:  CreatedEventMsg,
			NotifyAt: time.Now().String(),
			EventID:  int32(eventID),
		},
	}

	for _, id := range idsResp.IDs {
		req.UserIDs = append(req.UserIDs, int32(id))
	}

	_, err = h.NotificationService.CreateNotifications(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
