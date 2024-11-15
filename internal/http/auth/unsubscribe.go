package auth

import (
	"net/http"
	"strconv"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"

	"github.com/gorilla/mux"
)

// @Summary Отписаться от  пользователя
// @Description Отписаться от пользователя пользователя
// @Tags auth
// @Produce  json
// @Success 200
// @Failure 401 {object} httpErrors.HttpError "Unauthorized"
// @Failure 404 {object} httpErrors.HttpError "Not found"
// @Failure 409 {object} httpErrors.HttpError "Self subscription"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /users/subscribe/{id} [post]
func (h *AuthHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
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

	subscription := models.Subscription{
		SubscriberID: session.UserID,
		FollowsID:    id,
	}

	err = h.service.Unsubscribe(r.Context(), subscription)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrSubscriptionNotFound)
			return
		default:
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
