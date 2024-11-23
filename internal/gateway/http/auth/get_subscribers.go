package auth

import (
	"net/http"
	"strconv"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"

	"github.com/gorilla/mux"
)

// @Summary Подписки пользователя
// @Description Подписки пользователя
// @Tags auth
// @Produce  json
// @Success 200
// @Failure 404 {object} httpErrors.HttpError "Invalid ID"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /profile/subscribe/{id} [get]
func (h *AuthHandler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	paginationParams := utils.GetPaginationParams(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrInvalidID)
		return
	}

	users, err := h.service.GetSubscriptions(r.Context(), id)
	if err != nil {
		switch err {
		default:
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
			return
		}
	}

	resp := writeUsersResponse(users, paginationParams.Limit)

	utils.WriteResponse(w, http.StatusOK, resp)
}
