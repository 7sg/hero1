package database

import (
	"hero1/go-services/historical-events/api/proto/packages/eventpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Event is persisted in db
type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	Email       string             `bson:"email"`
	Environment string             `bson:"environment"`
	Component   string             `bson:"component"`
	Message     string             `bson:"message"`
	Data        string             `bson:"data"`
}

// SearchFilter for filtering events from db
type SearchFilter struct {
	Email       string          `bson:"email,omitempty"`
	Environment string          `bson:"environment,omitempty"`
	Component   string          `bson:"component,omitempty"`
	Text        primitive.Regex `bson:"message,omitempty"`
	Date        *DateFilter     `bson:"created_at,omitempty"`
}

// DateFilter is part of search filter
type DateFilter struct {
	DayStartTime time.Time `bson:"$gte,omitempty"`
	DayEndTime   time.Time `bson:"$lte,omitempty"`
}

// FromProto converts the proto search request to db search request
func (s *SearchFilter) FromProto(request *eventpb.SearchEventRequest) {
	s.Email = request.Email
	s.Environment = request.Environment
	s.Component = request.Component
}

// ToProto converts the db event to proto event
func (e *Event) ToProto() *eventpb.Event {
	return &eventpb.Event{
		Id:          e.ID.Hex(),
		CreatedAt:   e.CreatedAt.Unix(),
		Email:       e.Email,
		Environment: e.Environment,
		Component:   e.Component,
		Message:     e.Message,
		Data:        e.Data,
	}
}

// FromProto converts the proto event to db event
func (e *Event) FromProto(event *eventpb.Event) {
	e.CreatedAt = time.Unix(event.CreatedAt, 0)
	e.Email = event.Email
	e.Environment = event.Environment
	e.Component = event.Component
	e.Message = event.Message
	e.Data = event.Data
}
