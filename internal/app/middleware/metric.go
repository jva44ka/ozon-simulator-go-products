package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Metrics interface {
	ReportRequestInfo(methodName string, code string, duration time.Duration)
}

func ResponseTime(rm Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		code := status.Code(err).String()
		rm.ReportRequestInfo(info.FullMethod, code, duration)

		return resp, err
	}
}
