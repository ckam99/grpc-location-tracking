package interceptor

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcUnaryLogger(logger *slog.Logger) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		result, err := handler(ctx, req)
		duration := time.Since(time.Now())
		code := codes.Unknown
		if st, ok := status.FromError(err); ok {
			code = st.Code()
		}
		log := logger.With("op", "GrpcUnaryLogger", "method", info.FullMethod, "code", code, "duration", duration)

		if err != nil {
			log.Error(err.Error())
		} else {
			log.Info("received a gRPC request")
		}
		if code == codes.Internal {
			err = errors.New("internal error")
		}
		return result, err
	}

}
