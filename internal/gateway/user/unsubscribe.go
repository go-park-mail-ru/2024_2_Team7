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

func (h *UserHandlers) Unsubscribe(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.UserService.Unsubscribe(r.Context(), &subscription)
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
