//go:generate easyjson get_notifications.go
package events

import (
	"net/http"
	"time"

	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"
	"kudago/internal/models"
	pb "kudago/internal/notification/api"

	easyjson "github.com/mailru/easyjson"
)

//easyjson:json
type GetNotificationsResponse struct {
	Notifications []NotificationWithEvent `json:"notifications"`
}

//easyjson:json
type NotificationWithEvent struct {
	Notification models.Notification `json:"notification"`
	Event        models.Event        `json:"event"`
}

// @Summary Получение уведомлений по ID пользователя
// @Description Возвращает уведомления по идентификатору пользователя
// @Tags notifications
// @Produce  json
// @Success 200 {object} GetNotificationsResponse
// @Failure 404 {object} httpErrors.HttpError "Notification Not Found"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /notifications [get]
func (h EventHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	req := &pb.GetNotificationsRequest{UserID: int32(session.UserID)}
	notifications, err := h.NotificationService.GetNotifications(r.Context(), req)
	if err != nil {
		h.logger.Error(r.Context(), "get notifications", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	ids := make([]int, 0, len(notifications.Notifications))
	for _, n := range notifications.Notifications {
		ids = append(ids, int(n.EventID))
	}

	events, err := h.GetEventsByIDs(r.Context(), ids)
	if err != nil {
		h.logger.Error(r.Context(), "get events by ids", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	resp := writeNotificationsResponse(notifications.Notifications, events)
	w.WriteHeader(http.StatusOK)
	if _, err := easyjson.MarshalToWriter(&resp, w); err != nil {
		h.logger.Error(r.Context(), "get notifications", err)
	}
}

func writeNotificationsResponse(notifications []*pb.Notification, events map[int]models.Event) GetNotificationsResponse {
	layout := "2006-01-02 15:04:05 -0700 MST"

	response := GetNotificationsResponse{
		Notifications: make([]NotificationWithEvent, 0, len(notifications)),
	}

	for _, n := range notifications {
		notifyAt, _ := time.Parse(layout, n.NotifyAt)
		event, ok := events[int(n.EventID)]
		if !ok {
			continue
		}

		response.Notifications = append(response.Notifications, NotificationWithEvent{
			Notification: models.Notification{
				ID:       int(n.Id),
				UserID:   int(n.UserID),
				EventID:  int(n.EventID),
				Message:  n.Message,
				NotifyAt: notifyAt,
			},
			Event: event,
		})
	}

	return response
}
