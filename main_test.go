package telemetry

import (
	"context"
	"testing"

	"google.golang.org/grpc"
)

// TestInterceptor calls telemetry.Interceptor with an appropiate UnaryCall.
func TestInterceptor(t *testing.T) {
	ctx := context.Background()
	req := map[string]bool{"value": true}
	info := grpc.UnaryServerInfo{FullMethod: "/method"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return true, nil
	}

	res, err := Interceptor(ctx, req, &info, handler)

	if !res.(bool) {
		t.Fail()
	}

	if err != nil {
		t.Fatal(err.Error())
	}
}
