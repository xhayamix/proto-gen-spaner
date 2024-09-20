package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/xhayamix/proto-gen-spanner/cmd/di"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/entity"
	clientapi "github.com/xhayamix/proto-gen-spanner/pkg/domain/proto/client/api"
)

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	// configのセット
	ctx := context.Background()
	entity.InitConfig(ctx)

	app, err := di.InitializeApp(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	defer func() {
		app.SpannerDB.Close()
	}()

	s := grpc.NewServer()
	clientapi.RegisterUserServer(s, app.UserHandler)

	reflection.Register(s)

	// alb set path = "/grpc.health.v1.Health/Check?service=grpc-api"
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(s, healthSrv)
	healthSrv.SetServingStatus("grpc-api", healthpb.HealthCheckResponse_SERVING)
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
