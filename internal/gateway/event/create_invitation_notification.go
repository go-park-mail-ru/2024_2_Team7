package events

import (
	"net/http"
	"time"

	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"
	pb "kudago/internal/notification/api"

	"github.com/asaskevich/govalidator"
	"github.com/mailru/easyjson"
)

// @Summary Создание уведомления
// @Description Создание уведомления
// @Tags notifications
// @Accept  json
// @Param json body CreateNotificationRequest true "Данные для создания уведомления"
// @Success 200 {object} string "Notification created successfully"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /notifications [post]
func (h EventHandler) CreateInvitationNotification(w http.ResponseWriter, r *http.Request) {
	_, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	req := InviteNotificationRequest{}
	err := easyjson.UnmarshalFromReader(r.Body, &req)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}

	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		utils.ProcessValidationErrors(w, err)
		return
	}

	reqPB := &pb.CreateNotificationsRequest{
		UserIDs: []int32{int32(req.UserID)},
		Notification: &pb.Notification{
			Message:  InvitationMsg,
			NotifyAt: time.Now().String(),
			EventID:  int32(req.EventID),
		},
	}

	_, err = h.NotificationService.CreateNotifications(r.Context(), reqPB)
	if err != nil {
		h.logger.Error(r.Context(), "create notification", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	w.WriteHeader(http.StatusOK)
}
