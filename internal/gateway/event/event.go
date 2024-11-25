package events

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	pbEvent "kudago/internal/event/api"
	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"
	pbImage "kudago/internal/image/api"
	"kudago/internal/logger"

	"kudago/internal/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultPage   = 0
	defaultLimit  = 30
	maxUploadSize = 10 * 1024 * 1024 // 10Mb
)

type EventHandler struct {
	EventService pbEvent.EventServiceClient
	ImageService pbImage.ImageServiceClient
	logger       *logger.Logger
}

func NewEventHandlers(eventServiceAddr string, imageServiceAddr string, logger *logger.Logger) (*EventHandler, error) {
	eventConn, err := grpc.NewClient(eventServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	imageConn, err := grpc.NewClient(imageServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &EventHandler{
		EventService: pbEvent.NewEventServiceClient(eventConn),
		ImageService: pbImage.NewImageServiceClient(imageConn),
		logger:       logger,
	}, nil
}

var maxDate = time.Date(2030, 12, 31, 0, 0, 0, 0, time.UTC)

type EventResponse struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	Category    int      `json:"category_id"`
	Capacity    int      `json:"capacity"`
	Tag         []string `json:"tag"`
	AuthorID    int      `json:"author"`
	EventStart  string   `json:"event_start"`
	EventEnd    string   `json:"event_end"`
	ImageURL    string   `json:"image"`
	Latitude    float64  `json:"Latitude"`
	Longitude   float64  `json:"Longitude"`
}

type GetEventsResponse struct {
	Events []EventResponse `json:"events"`
}

type NewEventRequest struct {
	Title       string   `json:"title" valid:"required,length(3|100)"`
	Description string   `json:"description" valid:"required,length(3|1000)" `
	Location    string   `json:"location" valid:"length(3|100)"`
	Category    int      `json:"category_id" valid:"required,range(1|8)"`
	Capacity    int      `json:"capacity" valid:"range(1|20000)"`
	Tag         []string `json:"tag"`
	EventStart  string   `json:"event_start" valid:"rfc3339,required"`
	EventEnd    string   `json:"event_end" valid:"rfc3339,required"`
}

type NewEventResponse struct {
	Event EventResponse `json:"event"`
}

func checkNewEventRequest(req NewEventRequest) *httpErrors.HttpError {
	if len(req.Tag) > 3 {
		return httpErrors.ErrTooManyTags
	}

	for _, tag := range req.Tag {
		if len(tag) > 20 || len(tag) == 0 {
			return httpErrors.ErrBadTagLength
		}
	}

	eventStart, err := time.Parse(time.RFC3339, req.EventStart)
	if err != nil {
		return httpErrors.ErrInvalidTime
	}

	eventEnd, err := time.Parse(time.RFC3339, req.EventEnd)
	if err != nil {
		return httpErrors.ErrInvalidTime
	}

	if !eventEnd.After(eventStart) {
		return httpErrors.ErrEventStartAfterEventEnd
	}

	if eventStart.Before(time.Now()) || eventEnd.After(maxDate) {
		return httpErrors.ErrBadEventTiming
	}

	return nil
}

func parseEventData(r *http.Request) (NewEventRequest, *pbImage.UploadRequest, *httpErrors.HttpError) {
	var req NewEventRequest
	jsonData := r.FormValue("json")
	err := json.Unmarshal([]byte(jsonData), &req)
	if err != nil {
		return req, nil, httpErrors.ErrInvalidData
	}

	media, err := utils.HandleImageUpload(r)
	if err != nil {
		return req, nil, httpErrors.ErrInvalidImage
	}

	return req, media, nil
}

func toModelEvent(req NewEventRequest, authorID int) models.Event {
	return models.Event{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		EventStart:  req.EventStart,
		EventEnd:    req.EventEnd,
		AuthorID:    authorID,
		CategoryID:  req.Category,
		Capacity:    req.Capacity,
		Tag:         req.Tag,
	}
}

func toPBEvent(req NewEventRequest, authorID int) *pbEvent.Event {
	return &pbEvent.Event{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		EventStart:  req.EventStart,
		EventEnd:    req.EventEnd,
		AuthorID:    int32(authorID),
		CategoryID:  int32(req.Category),
		Capacity:    int32(req.Capacity),
		Tag:         req.Tag,
	}
}

func toEvent(event *pbEvent.Event) models.Event {
	return models.Event{
		Title:       event.Title,
		Description: event.Description,
		Location:    event.Location,
		EventStart:  event.EventStart,
		EventEnd:    event.EventEnd,
		AuthorID:    int(event.AuthorID),
		CategoryID:  int(event.CategoryID),
		Capacity:    int(event.Capacity),
		Tag:         event.Tag,
	}
}

func toIDPB(id int) *pbEvent.ID {
	return &pbEvent.ID{
		ID: int32(id),
	}
}

func eventToEventResponse(event *pbEvent.Event) EventResponse {
	return EventResponse{
		ID:          int(event.ID),
		Title:       event.Title,
		Description: event.Description,
		EventStart:  event.EventStart,
		EventEnd:    event.EventEnd,
		Tag:         event.Tag,
		AuthorID:    int(event.AuthorID),
		Category:    int(event.CategoryID),
		ImageURL:    event.Image,
		Capacity:    int(event.Capacity),
		Longitude:   float64(event.Latitude),
		Latitude:    float64(event.Latitude),
	}
}

func writeEventsResponse(events []*pbEvent.Event, limit int) GetEventsResponse {
	resp := GetEventsResponse{make([]EventResponse, 0, limit)}

	for _, event := range events {
		eventResp := eventToEventResponse(event)
		resp.Events = append(resp.Events, eventResp)
	}
	return resp
}

func GetQueryParamInt(r *http.Request, key string, defaultValue int) int {
	valueStr := r.URL.Query().Get(key)
	value, err := strconv.Atoi(valueStr)

	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}

func GetPaginationParams(r *http.Request) *pbEvent.PaginationParams {
	page := GetQueryParamInt(r, "page", defaultPage)
	limit := GetQueryParamInt(r, "limit", defaultLimit)
	offset := page * limit
	return &pbEvent.PaginationParams{
		Offset: int32(offset),
		Limit:  int32(limit),
	}
}

func (h *EventHandler) deleteImage(ctx context.Context, url string) {
	if url != "" {
		req := &pbImage.DeleteRequest{
			FileUrl: url,
		}
		h.ImageService.DeleteImage(ctx, req)
	}
}

func (h *EventHandler) uploadImage(ctx context.Context, media *pbImage.UploadRequest, w http.ResponseWriter) (string, error) {
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
