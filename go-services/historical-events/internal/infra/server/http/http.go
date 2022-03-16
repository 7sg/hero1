package http

import (
	"hero1/go-services/historical-events/api/proto/packages/eventpb"
	"hero1/go-services/historical-events/internal/config"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

type eventService interface {
	Save(context.Context, *eventpb.Event) (*eventpb.SaveEventResponse, error)
	Search(context.Context, *eventpb.SearchEventRequest) (*eventpb.SearchEventResponse, error)
}

// HTTP server struct,
type HTTP struct {
	config       *config.Config
	eventService eventService
}

type m struct {
	*runtime.JSONPb
}

// New create new http server
func New(config *config.Config, eventService eventService) *HTTP {
	return &HTTP{
		config:       config,
		eventService: eventService,
	}
}

// Start the http server
func (h *HTTP) Start() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Connect to the GRPC server
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%s", h.config.App.GRPCPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	defer func() {
		connCloseErr := conn.Close()
		if err != nil {
			log.Fatal(connCloseErr)
		}
	}()

	rmux := runtime.NewServeMux(
		runtime.WithMarshalerOption("*", &m{
			JSONPb: &runtime.JSONPb{
				EmitDefaults: true,
			},
		}),
	)

	eventServiceClient := eventpb.NewEventServiceClient(conn)
	err = eventpb.RegisterEventServiceHandlerClient(ctx, rmux, eventServiceClient)
	if err != nil {
		log.Printf("Error while registering client with grpc server, %v", err)
	}

	exposedMux := http.NewServeMux()

	swaggerUi := http.FileServer(http.Dir("go-services/historical-events/api/swagger-ui"))
	exposedMux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", swaggerUi))
	proto := http.FileServer(http.Dir("go-services/historical-events/api/proto"))
	exposedMux.Handle("/proto/", http.StripPrefix("/proto", proto))

	exposedMux.Handle("/", rmux)

	log.Printf("server starting to listen on %s", h.config.App.HTTPort)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", h.config.App.HTTPort), exposedMux)

	if err != nil {
		log.Fatal(err)
	}

}
