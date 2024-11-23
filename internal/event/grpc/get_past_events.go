package http

import (
	"context"

	pb "kudago/internal/event/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) GetPastEvents(ctx context.Context, req *pb.PaginationParams) (*pb.Events, error) {
	params := getPaginationParams(req)
	eventsData, err := s.getter.GetPastEvents(ctx, params)
	if err != nil {
		s.logger.Error(ctx, "get past events", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	event := writeEventsResponse(eventsData, params.Limit)

	return event, nil
}
