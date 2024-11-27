package grpc

import (
	"context"
	"errors"

	pb "kudago/internal/event/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) GetEventByID(ctx context.Context, req *pb.ID) (*pb.Event, error) {
	eventData, err := s.getter.GetEventByID(ctx, int(req.ID))
	if err != nil {
		if errors.Is(err, models.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, errEventNotFound)
		}
		s.logger.Error(ctx, "get event by id", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	event := eventToEventPB(eventData)

	return event, nil
}
