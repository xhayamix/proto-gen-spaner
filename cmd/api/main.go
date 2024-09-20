package main

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

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

	if checkSpannerConnection(app.SpannerDB) {
		log.Println("Connected to Spanner successfully.")
	} else {
		log.Println("Failed to connect to Spanner.")
	}

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

func checkSpannerConnection(client *spanner.Client) bool {
	ctx := context.Background()
	stmt := spanner.NewStatement("SELECT 1")
	ro := client.Single().WithTimestampBound(spanner.MaxStaleness(10 * time.Second))
	defer ro.Close()

	iter := ro.Query(ctx, stmt)
	defer iter.Stop()

	if row, err := iter.Next(); err != nil && err != iterator.Done {
		log.Println("Failed to query Spanner:", err)
		return false
	} else if row == nil {
		log.Println("No rows returned in query to Spanner")
		return false
	}
	return true
}
