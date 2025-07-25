package domain

import (
	"time"
)

// StoreEventMessage represents a message for store event processing through Kafka.
type StoreEventMessage struct {
	ProjectID  ProjectID            `json:"project_id"`
	EventID    EventID              `json:"event_id"`
	ReceivedAt time.Time            `json:"received_at"`
	EventData  map[string]any       `json:"event_data"`
	SourceIP   string               `json:"source_ip,omitempty"`
	UserAgent  string               `json:"user_agent,omitempty"`
	Processing StoreEventProcessing `json:"processing"`
}

// StoreEventProcessing contains processing information.
type StoreEventProcessing struct {
	Priority   int    `json:"priority"`    // Processing priority (1-10, where 10 is highest)
	WorkerPool string `json:"worker_pool"` // Worker pool for processing
}

// NewStoreEventMessage creates a new store event message.
func NewStoreEventMessage(projectID ProjectID, eventID EventID, eventData map[string]any) *StoreEventMessage {
	return &StoreEventMessage{
		ProjectID:  projectID,
		EventID:    eventID,
		ReceivedAt: time.Now(),
		EventData:  eventData,
		Processing: StoreEventProcessing{
			Priority:   5, // Default medium priority
			WorkerPool: "default",
		},
	}
}

// SetPriority sets processing priority.
func (sem *StoreEventMessage) SetPriority(priority int) {
	if priority < 1 {
		priority = 1
	}
	if priority > 10 {
		priority = 10
	}
	sem.Processing.Priority = priority
}

// SetWorkerPool sets worker pool.
func (sem *StoreEventMessage) SetWorkerPool(pool string) {
	sem.Processing.WorkerPool = pool
}
