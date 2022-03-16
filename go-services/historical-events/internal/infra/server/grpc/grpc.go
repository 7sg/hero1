package grpc

import (
	"hero1/go-services/historical-events/api/proto/packages/eventpb"
	"hero1/go-services/historical-events/internal/config"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type eventService interface {
	Save(context.Context, *eventpb.Event) (*eventpb.SaveEventResponse, error)
	Search(context.Context, *eventpb.SearchEventRequest) (*eventpb.SearchEventResponse, error)
}

// Grpc server struct,
type Grpc struct {
	config       *config.Config
	eventService eventService
}

// New create new grpc server struct
func New(config *config.Config, eventService eventService) *Grpc {
	return &Grpc{
		config:       config,
		eventService: eventService,
	}
}

// Start grpc server
func (g *Grpc) Start() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from err: %v", r)
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", g.config.App.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen on port %s %v", g.config.App.GRPCPort, err)
	}
	grpcServer := grpc.NewServer()
	eventpb.RegisterEventServiceServer(grpcServer, g.eventService)

	log.Printf("grpc server starting to listen on %s", g.config.App.GRPCPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc server over port %s %v", g.config.App.GRPCPort, err)
	}
}
