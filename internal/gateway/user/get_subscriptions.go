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

func (h *UserHandlers) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	paginationParams := utils.GetPaginationParams(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidID)
		return
	}

	users, err := h.UserService.GetSubscriptions(r.Context(), &pb.GetSubscriptionsRequest{ID: int32(id)})
	if err != nil {
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.NotFound:
				utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrSubscriptionNotFound)
				return
			}
		}

		h.logger.Error(r.Context(), "get subscriptions", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	resp := writeUsersResponse(users.Users, paginationParams.Limit)

	utils.WriteResponse(w, http.StatusOK, resp)
}
