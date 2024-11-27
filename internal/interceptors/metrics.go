package interceptors

import (
	"context"
	"time"

	"kudago/internal/metrics"

	"google.golang.org/grpc"
)

func MetricsUnaryInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)
		duration := time.Since(start).Seconds()

		statusCode := "OK"
		if err != nil {
			statusCode = "ERROR"
		}

		metrics.RequestDuration.WithLabelValues(info.FullMethod, "gRPC", serviceName, statusCode).Observe(duration)
		metrics.RequestCount.WithLabelValues(info.FullMethod, "gRPC", serviceName, statusCode).Inc()

		if err != nil {
			metrics.ErrorCount.WithLabelValues(info.FullMethod, "gRPC", serviceName, statusCode).Inc()
		}

		return resp, err
	}
}
