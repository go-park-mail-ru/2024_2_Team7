package handlers

import (
	"net/http"
	"strconv"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"
	pb "kudago/internal/user/api"

	"github.com/gorilla/mux"
)

// @Summary Подписка на пользователя
// @Description Подписка на пользователя
// @Tags auth
// @Produce  json
// @Success 200
// @Failure 401 {object} httpErrors.HttpError "Unauthorized"
// @Failure 404 {object} httpErrors.HttpError "Invalid ID"
// @Failure 409 {object} httpErrors.HttpError "Self subscription"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /users/subscribe/{id} [post]
func (h *UserHandlers) Subscribe(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusUnauthorized, httpErrors.ErrUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrInvalidID)
		return
	}

	if id == session.UserID {
		utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrSelfSubscription)
		return
	}

	subscription := pb.Subscription{
		SubscriberID: int32(session.UserID),
		FollowsID:    int32(id),
	}

	_, err = h.Gateway.UserService.Subscribe(r.Context(), &subscription)
	if err != nil {
		switch err {
		case models.ErrForeignKeyViolation:
			utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrUserNotFound)
			return
		case models.ErrNothingToInsert:
			utils.WriteResponse(w, http.StatusOK, httpErrors.ErrSubscriptionAlreadyExists)
			return
		default:
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
