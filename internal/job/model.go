package job

import (
	"encoding/json"
	"time"
)

type Status string
type Priority string

const (
	StatusPending    = "PENDING"
	StatusProcessing = "PROCESSING"
	StatusDelayed    = "DELAYED"
)

const (
	HighPriority   = "HIGH"
	MediumPriority = "MEDIUM"
	LowPriority    = "LOW"
)

type Job struct {
	ID           string          `json:"id" bson:"_id,omitempty"`
	Title        string          `json:"title" bson:"title"`
	Type         string          `json:"job_type" bson:"job_type"`
	JobStatus    Status          `json:"job_status" bson:"job_status"`
	Priority     Priority        `json:"priority" bson:"priority"`
	CurrAttempts int             `json:"curr_attempts" bson:"curr_attempts"`
	MaxAttempts  int             `json:"max_attempts" bson:"max_attempts"`
	Payload      json.RawMessage `json:"payload" bson:"payload"`
	CreatedAt    time.Time       `json:"created_at" bson:"created_at"`
}

type JobCreation struct {
	Title     string          `json:"title" bson:"title"`
	Type      string          `json:"job_type" bson:"job_type"`
	JobStatus Status          `json:"job_status" bson:"job_status"`
	Priority  Priority        `json:"priority" bson:"priority"`
	Payload   json.RawMessage `json:"payload" bson:"payload"`
}
