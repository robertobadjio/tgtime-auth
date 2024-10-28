package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/robertobadjio/tgtime-auth/internal/logger"
)

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	now := time.Now()

	res, err := handler(ctx, req)
	if err != nil {
		logger.Error("request", err.Error(), "method", info.FullMethod, "req", req)
	}

	logger.Info("request", "", "method", info.FullMethod, "req", req, "res", res, "duration", time.Since(now))

	return res, err
}
