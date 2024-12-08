package handlers

import (
	"log"
	"net/http"
	"time"

	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"
	"kudago/internal/models"
	pb "kudago/internal/notification/api"
)

type GetNotificationsResponse struct {
	Notifications []models.Notification `json:"notifications"`
}

// @Summary Получение уведомлений по ID пользователя
// @Description Возвращает уведомления по идентификатору пользователя
// @Tags notifications
// @Produce  json
// @Success 200 {object} GetNotificationsResponse 
// @Failure 404 {object} httpErrors.HttpError "Notification Not Found"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /notifications [get]
func (h NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
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
	log.Println(notifications)
	resp := writeNotificationsResponse(notifications.Notifications)
	utils.WriteResponse(w, http.StatusOK, resp)
}

func writeNotificationsResponse(notifications []*pb.Notification) GetNotificationsResponse {
	layout := "2006-01-02 15:04:05 -0700 MST" 
		response := GetNotificationsResponse{
		Notifications: make([]models.Notification, len(notifications)),
	}

	for i, n := range notifications {
		notifyAt, err := time.Parse(layout, n.NotifyAt)
		log.Println(err)
		response.Notifications[i] = models.Notification{
			ID:       int(n.Id),
			UserID:   int(n.UserID),
			EventID:  int(n.EventID),
			Message:  n.Message,
			NotifyAt: notifyAt,
		}
	}

	return response
}
