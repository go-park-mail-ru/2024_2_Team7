package handlers

import (
	"context"
	"net/http"

	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"
	pbImage "kudago/internal/image/api"
	"kudago/internal/models"
	pb "kudago/internal/user/api"

	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"

	"github.com/asaskevich/govalidator"
)

func (h *UserHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	req, media, reqErr := parseUpdateData(r)
	if reqErr != nil {
		utils.WriteResponse(w, http.StatusBadRequest, reqErr)
		return
	}

	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		utils.ProcessValidationErrors(w, err)
		return
	}

	url, err := h.uploadImage(r.Context(), media, w)
	if err != nil {
		return
	}

	req.AvatarUrl = url
	req.ID = int32(session.UserID)

	user, err := h.UserService.UpdateUser(r.Context(), req)
	if err != nil {
		h.deleteImage(r.Context(), url)
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.AlreadyExists:
				utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrUsernameIsAlredyTaken)
				return
			case grpcCodes.Internal:
				h.logger.Error(r.Context(), "update user", st.Err())
				utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
				return
			}
			h.logger.Error(r.Context(), "update user", err)
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
			return
		}
	}

	resp := userToUserResponse(user)

	utils.WriteResponse(w, http.StatusOK, resp)
	return
}

func parseUpdateData(r *http.Request) (*pb.User, *pbImage.UploadRequest, *httpErrors.HttpError) {
	var req models.User
	jsonData := r.FormValue("json")

	err := req.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		return nil, nil, httpErrors.ErrInvalidData
	}

	media, err := utils.HandleImageUpload(r)
	if err != nil {
		return nil, nil, httpErrors.ErrInvalidImage
	}

	user := &pb.User{
		ID:        int32(req.ID),
		Username:  req.Username,
		Email:     req.Email,
		AvatarUrl: req.ImageURL,
	}
	return user, media, nil
}

func (h *UserHandlers) uploadImage(ctx context.Context, media *pbImage.UploadRequest, w http.ResponseWriter) (string, error) {
	if media != nil {
		url, err := h.ImageService.UploadImage(ctx, media)
		if err != nil {
			switch err {
			case models.ErrInvalidImage:
				utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidImage)
			case models.ErrInvalidImageFormat:
				utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidImageFormat)
			default:
				h.logger.Error(ctx, "upload image", err)
				utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
			}
			return "", err
		}
		return url.FileUrl, nil
	}
	return "", nil
}

func (h *UserHandlers) deleteImage(ctx context.Context, url string) {
	if url != "" {
		req := &pbImage.DeleteRequest{
			FileUrl: url,
		}
		h.ImageService.DeleteImage(ctx, req)
	}
}
