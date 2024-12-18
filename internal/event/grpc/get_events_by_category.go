package grpc

import (
	"context"

	pb "kudago/internal/event/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) GetEventsByCategory(ctx context.Context, req *pb.GetEventsByCategoryRequest) (*pb.Events, error) {
	params := getPaginationParams(req.Params)
	eventsData, err := s.getter.GetEventsByCategory(ctx, int(req.CategoryID), params)
	if err != nil {
		s.logger.Error(ctx, "get events by category", err)
		return nil, status.Error(codes.Internal, ErrInternal)
	}

	event := writeEventsResponse(eventsData, params.Limit)

	return event, nil
}
