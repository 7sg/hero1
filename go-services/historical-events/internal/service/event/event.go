package event

import (
	"hero1/go-services/historical-events/api/proto/packages/eventpb"
	"hero1/go-services/historical-events/internal/domain/database"
	"context"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
	"google.golang.org/grpc/metadata"
)

// Service for events
type Service struct {
	eventRepository eventRepository
}

// New will Create and instance of event.Service
func New(repository eventRepository) *Service {
	return &Service{eventRepository: repository}
}

type eventRepository interface {
	Save(context.Context, *database.Event) error
	Get(context.Context, *database.SearchFilter) ([]*database.Event, error)
}

// Save service is responsible for saving event
func (s *Service) Save(ctx context.Context, event *eventpb.Event) (*eventpb.SaveEventResponse, error) {
	dbEvent := &database.Event{}
	dbEvent.FromProto(event)
	err := s.eventRepository.Save(ctx, dbEvent)
	if err != nil {
		log.Printf("error from eventRepository on Save, err:%+v", err)
		return nil, err
	}
	return &eventpb.SaveEventResponse{}, err
}

// Search service is responsible for searching events based on filters
func (s *Service) Search(ctx context.Context, searchRequest *eventpb.SearchEventRequest) (*eventpb.SearchEventResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	log.Printf("metadat is %+v, ok is %v", md, ok)
	searchFilter := &database.SearchFilter{}
	searchFilter.FromProto(searchRequest)
	searchFilter.Text = primitive.Regex{Pattern: searchRequest.Text, Options: ""}

	if searchRequest.GetDate() != nil {
		searchDate, _ := ptypes.Timestamp(searchRequest.GetDate())
		dayStartTime := time.Date(searchDate.Year(), searchDate.Month(), searchDate.Day(), 0, 0, 0, 0, time.UTC)
		dayEndTime := time.Date(searchDate.Year(), searchDate.Month(), searchDate.Day(), 23, 59, 59, 999999999, time.UTC)

		searchFilter.Date = &database.DateFilter{
			DayStartTime: dayStartTime,
			DayEndTime:   dayEndTime,
		}
	}

	events, err := s.eventRepository.Get(ctx, searchFilter)
	if err != nil {
		log.Printf("error from eventRepository on Get, err:%+v", err)
		return nil, err
	}

	var protoEvents []*eventpb.Event
	for _, event := range events {
		protoEvents = append(protoEvents, event.ToProto())
	}

	return &eventpb.SearchEventResponse{Events: protoEvents}, err
}

/* In case business logic become complex
   - If it is related to event(single responsibility), then i will add functions in this package,
   - otherwise i will create new services and repository, which will have its own packages.
*/

/* Future business logic addition will grow domain, service and repository layers.
   I will identify new storage level entities, and create respective domain , repositories and services in respective packages.
*/
