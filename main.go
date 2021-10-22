package telemetry

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const MethodHealthCheck = "/grpc.health.v1.Health/Check"

func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == MethodHealthCheck {
		return handler(ctx, req)
	}

	// Start time
	start := time.Now()

	// Calls the handler
	res, err := handler(ctx, req)

	lat := float64(time.Since(start).Nanoseconds()) / float64(1e6)
	md, _ := metadata.FromIncomingContext(ctx)
	code := status.Convert(err).Code()

	// Log request info
	var logger = func(event *zerolog.Event) *zerolog.Event {
		return event.Str("status", code.String()).
			Float64("latency", lat).
			Interface("request", req).
			Interface("metadata", md)
	}

	if err != nil {
		switch code {
		case codes.InvalidArgument, codes.Unauthenticated, codes.PermissionDenied, codes.NotFound:
			logger(log.Warn()).Err(err).Msg(info.FullMethod)

		default:
			logger(log.Error()).Err(err).Msg(info.FullMethod)
		}

		return nil, err
	}

	logger(log.Info()).Msg(info.FullMethod)

	return res, nil
}
