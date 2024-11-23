package http

import (
	"context"

	pb "kudago/internal/event/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) GetEventByID(ctx context.Context, req *pb.ID) (*pb.Event, error) {
	eventData, err := s.getter.GetEventByID(ctx, int(req.ID))
	if err != nil {
		s.logger.Error(ctx, "get event by id", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	event := eventToEventPB(eventData)

	return event, nil
}
