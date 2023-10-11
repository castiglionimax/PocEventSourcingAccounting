package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Model struct {
	EventID     string    `json:"event_id" bson:"event_id"`
	EventType   string    `json:"event_type" bson:"event_type"`
	AggregateID string    `json:"aggregate_id" bson:"aggregate_id"`
	Time        time.Time `json:"time" bson:"time"`
	Data        any       `json:"data" bson:"data"`
}

func NewModel(eventType, aggregateID string, data any) Model {
	uniqueID := uuid.New().String()
	timestamp := time.Now()
	strParts := strings.Split(uniqueID, "-")

	return Model{
		EventID:     fmt.Sprintf("%s-%d", strParts[0], timestamp.UnixNano()/int64(time.Millisecond)),
		EventType:   eventType,
		AggregateID: aggregateID,
		Time:        timestamp,
		Data:        data,
	}
}
