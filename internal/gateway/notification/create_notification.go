package handlers

import (
	"encoding/json"
	"net/http"

	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"
	pb "kudago/internal/notification/api"

	"github.com/asaskevich/govalidator"
)

type CreateNotificationRequest struct {
	UserID   int    `json:"user_id" valid:"range(1|20000)"`
	EventID  int    `json:"event_id" valid:"range(1|20000)"`
	Message  string `json:"message" valid:"required,length(3|1000)"`
	NotifyAt string `json:"notify_at" valid:"rfc3339,required"`
}

// @Summary Создание уведомления
// @Description Создание уведомления
// @Tags notifications
// @Accept  json
// @Param json body CreateNotificationRequest true "Данные для создания уведомления"
// @Success 200 {object} string "Notification created successfully"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /notifications [post]
func (h NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	_, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	var req CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}

	_, err := govalidator.ValidateStruct(&req)
	if err != nil {
		utils.ProcessValidationErrors(w, err)
		return
	}

	reqPB := &pb.CreateNotificationRequest{
		Notification: &pb.Notification{
			UserID:   int32(req.UserID),
			EventID:  int32(req.EventID),
			Message:  req.Message,
			NotifyAt: req.NotifyAt,
		},
	}

	_, err = h.NotificationService.CreateNotification(r.Context(), reqPB)
	if err != nil {
		h.logger.Error(r.Context(), "create notification", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	w.WriteHeader(http.StatusOK)
}
