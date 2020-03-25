package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/phisuite/data.gw"
	"github.com/phisuite/schema.gw"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

type registerHandler func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
type registerService func(context.Context, *runtime.ServeMux)

const GrpcPort = 50051

func main() {
	port := flag.Int("port", 80, "the server port")
	flag.Parse()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	router := runtime.NewServeMux()

	registrants := []registerService{
		build("entity-inspector", []registerHandler{
			data.RegisterEntityReadAPIHandler,
		}),
		build("event-publisher", []registerHandler{
			data.RegisterEventAPIHandler,
		}),
		build("schema-inspector", []registerHandler{
			schema.RegisterEntityReadAPIHandler,
			schema.RegisterEventReadAPIHandler,
		}),
	}

	for _, register := range registrants {
		register(ctx, router)
	}

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router)
	if err != nil {
		log.Fatalf("failed to serve http: %v", err)
	}
}

func build(service string, registrants []registerHandler) registerService {
	return func(ctx context.Context, router *runtime.ServeMux) {
		serviceAddr := fmt.Sprintf("%s:%d", service, GrpcPort)
		conn, err := grpc.DialContext(ctx, serviceAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("failed to establish connection with %s: %v", serviceAddr, err)
		}

		for _, register := range registrants {
			err = register(ctx, router, conn)
			if err != nil {
				log.Printf("failed to register %s: %v", serviceAddr, err)
			}
		}
	}
}
