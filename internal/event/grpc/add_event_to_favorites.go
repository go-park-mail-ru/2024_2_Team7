package grpc

import (
	"context"

	pb "kudago/internal/event/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) AddEventToFavorites(ctx context.Context, req *pb.FavoriteEvent) (*pb.Empty, error) {
	newFavorite := favoritePBToFavorite(req)

	err := s.service.AddEventToFavorites(ctx, newFavorite)
	if err != nil {
		if err != nil {
			s.logger.Error(ctx, "add event to favorites", err)
			switch err {
			case models.ErrForeignKeyViolation:
				return nil, status.Error(codes.NotFound, ErrEventNotFound)
			case models.ErrNothingToInsert:
				return nil, status.Error(codes.AlreadyExists, ErrAlreadyInFavorites)
			default:
				return nil, status.Error(codes.Internal, ErrInternal)
			}
		}
	}

	return nil, nil
}
