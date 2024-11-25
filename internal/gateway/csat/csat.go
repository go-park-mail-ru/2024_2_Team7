package csat

import (
	"encoding/json"
	"net/http"

	pb "kudago/internal/csat/api"
	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"
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
			case grpcCodes.Internal:
				h.logger.Error(r.Context(), "get test", st.Err())
				utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
				return
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

func (h *CSATHandlers) GetStatistics(w http.ResponseWriter, r *http.Request) {
	_, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusOK, httpErrors.ErrUnauthorized)
		return
	}

	stats, err := h.CSATService.GetStatistics(r.Context(), nil)
	if err != nil {
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.Internal:
				h.logger.Error(r.Context(), "get stats", st.Err())
				utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
				return
			default:
				h.logger.Error(r.Context(), "get stats", st.Err())
				utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
				return
			}
		}
		h.logger.Error(r.Context(), "get stats", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	resp := toStatisticsModel(stats)

	utils.WriteResponse(w, http.StatusOK, resp)
	return
}

func (h *CSATHandlers) AddAnswers(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusOK, httpErrors.ErrUnauthorized)
		return
	}

	var answers models.AddAnswers
	err := json.NewDecoder(r.Body).Decode(&answers)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrUserAlreadyDidTest)
		return
	}

	answersReq := make([]*pb.Answer, 0, len(answers.Answers))
	for _, answer := range answers.Answers {
		temp := &pb.Answer{
			QuestionID: int32(answer.QuestionID),
			Value:      int32(answer.Value),
		}
		answersReq = append(answersReq, temp)
	}

	req := &pb.AddAnswersRequest{
		UserID:  int32(session.UserID),
		Answers: answersReq,
	}

	_, err = h.CSATService.AddAnswers(r.Context(), req)
	if err != nil {
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.AlreadyExists:
				h.logger.Error(r.Context(), "get stats", st.Err())
				utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrUserAlreadyDidTest)
				return
			case grpcCodes.Internal:
				h.logger.Error(r.Context(), "get test", st.Err())
				utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
				return
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

	w.WriteHeader(http.StatusOK)
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

func toStatisticsModel(stats *pb.GetStatisticsResponse) models.Statistics {
	statistics := make([]models.Stats, 0, len(stats.Statistics))

	for _, stat := range stats.Statistics {
		q := toStatsModel(stat)
		statistics = append(statistics, q)
	}

	return models.Statistics{
		Statistics: statistics,
	}
}

func toQuestionModel(question *pb.Question) models.Question {
	return models.Question{
		ID:   int(question.Id),
		Text: question.Text,
	}
}

func toStatsModel(stats *pb.Stats) models.Stats {
	return models.Stats{
		ID:       int(stats.ID),
		Question: stats.Question,
		Value:    int(stats.Value),
	}
}
