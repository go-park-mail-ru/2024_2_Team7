package handlers

import (
	"net/http"
	"strconv"

	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"
	pb "kudago/internal/user/api"

	"github.com/gorilla/mux"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
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
// @Router /profile/subscribe/{id} [post]
func (h *UserHandlers) Subscribe(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusUnauthorized, httpErrors.ErrUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidID)
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

	_, err = h.UserService.Subscribe(r.Context(), &subscription)
	if err != nil {
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.NotFound:
				utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrUserNotFound)
				return
			case grpcCodes.AlreadyExists:
				utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrSubscriptionAlreadyExists)
				return
			case grpcCodes.Internal:
				h.logger.Error(r.Context(), "subscribe", st.Err())
				utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
				return
			}
		}

		h.logger.Error(r.Context(), "subscribe", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	w.WriteHeader(http.StatusOK)
}
