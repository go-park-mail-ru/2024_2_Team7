package http

import (
	"context"

	pb "kudago/internal/event/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) UpdateEvent(ctx context.Context, req *pb.Event) (*pb.Event, error) {
	newEvent := eventPBToEvent(req)

	eventData, err := s.service.UpdateEvent(ctx, newEvent)
	if err != nil {
		s.logger.Error(ctx, "update event", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	event := eventToEventPB(eventData)

	return event, nil
}
