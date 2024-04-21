package suite

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	ssov1 "orchestrator/protos/gen/go/sso"
	"orchestrator/tracing"
	"testing"
	"time"
)

const (
	grpcHost = "localhost"
)

type Suite struct {
	*testing.T                  // потребуется для вызова методов *testing.T внутри Suite
	AuthClient ssov1.AuthClient // Клиент для взаимодействия с grpc_transport - сервером
}

var tracer = otel.Tracer("testing client")

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()   // in case of failure in one of tests to form rightly stack trace
	t.Parallel() //can run test in parallel to increase performance

	tp, err := tracing.Init("testing client")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)

	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"client-id", "web-api-client-us-east-1",
		"user-id", "some-test-user-id",
	)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	_, span := tracer.Start(ctx, "Testing Client")
	defer span.End()
	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

	// context will be canceled when tests are stopped
	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	// create grpc_transport client
	cc, err := grpc.DialContext(
		context.Background(),
		grpcAddress(),
		//use insecure connection during test
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		t.Fatalf("grpc_transport server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAddress() string {
	return net.JoinHostPort("localhost", "44044")
}
