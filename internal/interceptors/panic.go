package interceptors

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func PanicRecoveryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic intercepted in method %s: %v", info.FullMethod, r)
		}
	}()

	return handler(ctx, req)
}
