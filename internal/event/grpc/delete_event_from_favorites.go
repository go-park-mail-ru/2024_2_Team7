package http

import (
	"context"
	"errors"

	pb "kudago/internal/event/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) DeleteEventFromFavorites(ctx context.Context, req *pb.FavoriteEvent) (*pb.Empty, error) {
	newFavorite := favoritePBToFavorite(req)

	err := s.service.DeleteEventFromFavorites(ctx, newFavorite)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, errEventNotFound)
		}
		s.logger.Error(ctx, "delete event from favorites", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	return nil, nil
}
