package main

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
	"time"
)

func grpcAddress() string {
	return net.JoinHostPort("localhost", "44044")
}

func main() {
	var tracer = otel.Tracer("testing client")
	tp, err := tracing.Init("testing client")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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

	// create grpc_transport client
	cc, err := grpc.DialContext(
		context.Background(),
		grpcAddress(),
		//use insecure connection during test
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	authClient := ssov1.NewAuthClient(cc)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Imx1bmFyYXRoQHJpcHBpbi5uZXQiLCJleHAiOjE3MTM0Njg2MjYsInRva2VuX3R5cGUiOiJhY2Nlc3MiLCJ1aWQiOjI0fQ.4Z07pCvO6ohZLm6NCJhSoofW453YUHzcSeDwTdQTkb4"
	respIsValid, err := authClient.Validate(ctx, &ssov1.ValidateRequest{Token: token})
	success := respIsValid.GetSuccess()
	fmt.Println(success, err)

}
