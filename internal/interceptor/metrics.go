package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/robertobadjio/tgtime-auth/internal/metric"
)

// MetricsInterceptor ???
func MetricsInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	metric.IncRequestCounter()

	timeStart := time.Now()

	res, err := handler(ctx, req)
	diffTime := time.Since(timeStart)

	if err != nil {
		metric.IncResponseCounter("error", info.FullMethod)
		metric.HistogramResponseTimeObserve("error", diffTime.Seconds())
	} else {
		metric.IncResponseCounter("success", info.FullMethod)
		metric.HistogramResponseTimeObserve("success", diffTime.Seconds())
	}

	return res, err
}
