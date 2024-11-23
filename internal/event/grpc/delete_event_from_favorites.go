package http

import (
	"context"

	pb "kudago/internal/event/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) DeleteEventFromFavorites(ctx context.Context, req *pb.FavoriteEvent) (*pb.Empty, error) {
	newFavorite := favoritePBToFavorite(req)

	err := s.service.DeleteEventFromFavorites(ctx, newFavorite)
	if err != nil {
		s.logger.Error(ctx, "delete event to favorites", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	return nil, nil
}
