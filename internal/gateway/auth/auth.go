package handlers

import (
	"context"
	"net/http"
	"regexp"
	"time"

	auth "kudago/internal/auth/api"
	pb "kudago/internal/auth/api"
	"kudago/internal/gateway/utils"
	httpErrors "kudago/internal/http/errors"
	pbImage "kudago/internal/image/api"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthHandlers struct {
	AuthService  pb.AuthServiceClient
	ImageService pbImage.ImageServiceClient
	logger       *logger.Logger
}

var validPasswordRegex = regexp.MustCompile(`^[a-zA-Z0-9+\-*/.;=\]\[\}\{\?]+$`)

func init() {
	govalidator.TagMap["password"] = govalidator.Validator(func(str string) bool {
		return validPasswordRegex.MatchString(str)
	})
}

func NewAuthHandlers(authServiceAddr string, imageServiceAddr string, logger *logger.Logger) (*AuthHandlers, error) {
	authConn, err := grpc.NewClient(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	imageConn, err := grpc.NewClient(imageServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &AuthHandlers{
		AuthService:  auth.NewAuthServiceClient(authConn),
		ImageService: pbImage.NewImageServiceClient(imageConn),
		logger:       logger,
	}, nil
}

type AuthResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	ImageURL string `json:"image"`
}

func userToUserResponse(user *pb.User) AuthResponse {
	resp := AuthResponse{
		User: UserResponse{
			ID:       int(user.ID),
			Username: user.Username,
			Email:    user.Email,
			ImageURL: user.AvatarUrl,
		},
	}
	return resp
}

func (h *AuthHandlers) setSessionCookie(w http.ResponseWriter, r *http.Request, ID int) error {
	req := &pb.CreateSessionRequest{ID: int32(ID)}

	session, err := h.AuthService.CreateSession(r.Context(), req)
	if err != nil {
		return models.ErrInternal
	}

	expires, err := time.Parse(time.RFC3339, session.Expires)
	if err != nil {
		return models.ErrInternal
	}

	http.SetCookie(w, &http.Cookie{
		Name:     models.SessionToken,
		Value:    session.Token,
		Expires:  expires,
		HttpOnly: true,
	})
	return nil
}

func (h *AuthHandlers) deleteImage(ctx context.Context, url string) {
	if url != "" {
		req := &pbImage.DeleteRequest{
			FileUrl: url,
		}
		h.ImageService.DeleteImage(ctx, req)
	}
}

func (h *AuthHandlers) uploadImage(ctx context.Context, media *pbImage.UploadRequest, w http.ResponseWriter) (string, error) {
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
