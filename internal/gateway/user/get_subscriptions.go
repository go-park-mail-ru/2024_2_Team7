package handlers

import (
	"net/http"
	"strconv"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	pb "kudago/internal/user/api"

	"github.com/gorilla/mux"
)

func (h *UserHandlers) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	paginationParams := utils.GetPaginationParams(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrInvalidID)
		return
	}

	users, err := h.Gateway.UserService.GetSubscriptions(r.Context(), &pb.GetSubscriptionsRequest{ID: int32(id)})
	if err != nil {
		switch err {
		default:
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
			return
		}
	}

	resp := writeUsersResponse(users.Users, paginationParams.Limit)

	utils.WriteResponse(w, http.StatusOK, resp)
}
