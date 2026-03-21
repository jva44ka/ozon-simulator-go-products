package middleware

import (
	"context"
	"log/slog"

	"github.com/jva44ka/ozon-simulator-go-products/internal/infra/config"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func Logger(cfg *config.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if cfg.Logging.LogRequestBody {
			body, _ := protojson.Marshal((req).(proto.Message))
			slog.InfoContext(ctx, "request", "method", info.FullMethod, "body", string(body))
		}

		if resp, err = handler(ctx, req); err != nil {
			slog.ErrorContext(ctx, "request failed", "method", info.FullMethod, "err", err)
		}

		if cfg.Logging.LogResponseBody {
			if v, ok := resp.(proto.Message); ok {
				body, _ := protojson.Marshal(v)
				slog.InfoContext(ctx, "response", "method", info.FullMethod, "body", string(body))
			}
		}

		return resp, err
	}
}
