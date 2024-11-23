package grpc

import (
	"bytes"
	"context"

	pb "kudago/internal/image/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) UploadImage(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	file := bytes.NewReader(req.File)

	mediaFile := models.MediaFile{
		Filename: req.Filename,
		File:     &readSeekCloser{Reader: file},
	}

	url, err := s.service.UploadImage(ctx, mediaFile)
	if err != nil {
		s.logger.Error(ctx, "upload image", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	resp := &pb.UploadResponse{
		FileUrl: url,
	}
	return resp, nil
}

type readSeekCloser struct {
	*bytes.Reader
}

func (r *readSeekCloser) Close() error {
	return nil
}
