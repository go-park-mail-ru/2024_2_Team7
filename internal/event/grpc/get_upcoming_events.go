package grpc

import (
	"context"

	pb "kudago/internal/event/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) GetUpcomingEvents(ctx context.Context, req *pb.PaginationParams) (*pb.Events, error) {
	params := getPaginationParams(req)
	eventsData, err := s.getter.GetUpcomingEvents(ctx, params)
	if err != nil {
		s.logger.Error(ctx, "get upcoming events", err)
		return nil, status.Error(codes.Internal, ErrInternal)
	}

	event := writeEventsResponse(eventsData, params.Limit)

	return event, nil
}
