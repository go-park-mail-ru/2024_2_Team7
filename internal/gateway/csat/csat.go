package csat

import (
	"net/http"

	pb "kudago/internal/csat/api"
	"kudago/internal/gateway/utils"
	httpErrors "kudago/internal/http/errors"
	"kudago/internal/logger"
	"kudago/internal/models"

	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	grpcStatus "google.golang.org/grpc/status"
)

type CSATHandlers struct {
	CSATService pb.CSATServiceClient
	logger      *logger.Logger
}

func NewCSATHandlers(csatServiceAddr string, logger *logger.Logger) (*CSATHandlers, error) {
	csatConn, err := grpc.NewClient(csatServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &CSATHandlers{
		CSATService: pb.NewCSATServiceClient(csatConn),
		logger:      logger,
	}, nil
}

func (h *CSATHandlers) GetTest(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusOK, httpErrors.ErrUnauthorized)
		return
	}

	query := r.URL.Query().Get("query")

	req := &pb.GetTestRequest{
		Query:  query,
		UserID: int32(session.UserID),
	}

	test, err := h.CSATService.GetTest(r.Context(), req)
	if err != nil {
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.NotFound:
				utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrTestNotFound)
				return
			// case grpcCodes.Internal:
			// 	h.logger.Error(r.Context(), "check session", st.Err())
			// 	utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
			// 	return
			default:
				h.logger.Error(r.Context(), "get test", st.Err())
				utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
				return
			}
		}
		h.logger.Error(r.Context(), "get test", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	resp := toTestModel(test)

	utils.WriteResponse(w, http.StatusOK, resp)
	return
}

func toTestModel(test *pb.GetTestResponse) models.Test {
	questions := make([]models.Question, 0, len(test.Questions))

	for _, question := range test.Questions {
		q := toQuestionModel(question)
		questions = append(questions, q)
	}

	return models.Test{
		ID:        int(test.Id),
		Title:     test.Title,
		Questions: questions,
	}
}

func toQuestionModel(question *pb.Question) models.Question {
	return models.Question{
		ID:   int(question.Id),
		Text: question.Text,
	}
}
