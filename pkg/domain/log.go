package audit

import (
	"time"
)

type LogItem struct {
	Action    string    `bson:"action" json:"action"`
	Entity    string    `bson:"entity" json:"entity"`
	EntityID  int64     `bson:"entity_id" json:"entity_id"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type LogMessage struct {
	Context map[string]interface{} `json:"context"`
	LogItem LogItem         `json:"log_item"`
}