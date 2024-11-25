package http

import (
	"context"

	pb "kudago/internal/event/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) SearchEvents(ctx context.Context, req *pb.SearchParams) (*pb.Events, error) {
	params := getPaginationParams(req.Params)

	searchParams := models.SearchParams{
		Query:        req.Query,
		EventStart:   req.EventStart,
		EventEnd:     req.EventEnd,
		Tags:         req.Tag,
		Category:     int(req.CategoryID),
		LatitudeMin:  float64(req.LatitudeMin),
		LatitudeMax:  float64(req.LatitudeMax),
		LongitudeMin: float64(req.LongitudeMin),
		LongitudeMax: float64(req.LongitudeMax),
	}

	eventsData, err := s.service.SearchEvents(ctx, searchParams, params)
	if err != nil {
		s.logger.Error(ctx, "search events", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	event := writeEventsResponse(eventsData, params.Limit)

	return event, nil
}
