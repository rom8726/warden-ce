package domain

import (
	"encoding/json"
	"time"
)

// EnvelopeMessage represents a message for envelope processing through Kafka.
type EnvelopeMessage struct {
	ProjectID  ProjectID          `json:"project_id"`
	EnvelopeID string             `json:"envelope_id"`
	ReceivedAt time.Time          `json:"received_at"`
	Data       []byte             `json:"data"`
	Header     map[string]any     `json:"header"`
	Items      []EnvelopeItem     `json:"items"`
	SourceIP   string             `json:"source_ip,omitempty"`
	UserAgent  string             `json:"user_agent,omitempty"`
	Processing EnvelopeProcessing `json:"processing"`
}

// EnvelopeItem represents an envelope item.
type EnvelopeItem struct {
	Type    string          `json:"type"`
	Length  int             `json:"length"`
	Payload json.RawMessage `json:"payload"`
}

// EnvelopeProcessing contains processing information.
type EnvelopeProcessing struct {
	Priority   int    `json:"priority"`    // Processing priority (1-10, where 10 is highest)
	RetryCount int    `json:"retry_count"` // Number of processing attempts
	MaxRetries int    `json:"max_retries"` // Maximum number of attempts
	Deadline   int64  `json:"deadline"`    // Deadline in Unix timestamp
	WorkerPool string `json:"worker_pool"` // Worker pool for processing
}

// NewEnvelopeMessage creates a new envelope message.
func NewEnvelopeMessage(projectID ProjectID, data []byte) *EnvelopeMessage {
	return &EnvelopeMessage{
		ProjectID:  projectID,
		EnvelopeID: generateEnvelopeID(),
		ReceivedAt: time.Now(),
		Data:       data,
		Processing: EnvelopeProcessing{
			Priority:   5, // Default medium priority
			RetryCount: 0,
			MaxRetries: 3,
			Deadline:   time.Now().Add(15 * time.Minute).Unix(), // 15 minutes to process
			WorkerPool: "default",
		},
	}
}

// SetPriority sets processing priority.
func (em *EnvelopeMessage) SetPriority(priority int) {
	if priority < 1 {
		priority = 1
	}
	if priority > 10 {
		priority = 10
	}
	em.Processing.Priority = priority
}

// SetWorkerPool sets worker pool.
func (em *EnvelopeMessage) SetWorkerPool(pool string) {
	em.Processing.WorkerPool = pool
}

// SetDeadline sets processing deadline.
func (em *EnvelopeMessage) SetDeadline(deadline time.Time) {
	em.Processing.Deadline = deadline.Unix()
}

// IsExpired checks if deadline has expired.
func (em *EnvelopeMessage) IsExpired() bool {
	return time.Now().Unix() > em.Processing.Deadline
}

// CanRetry checks if processing can be retried.
func (em *EnvelopeMessage) CanRetry() bool {
	return em.Processing.RetryCount < em.Processing.MaxRetries
}

// IncrementRetry increases retry counter.
func (em *EnvelopeMessage) IncrementRetry() {
	em.Processing.RetryCount++
}

// generateEnvelopeID generates unique ID for envelope.
func generateEnvelopeID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates random string of given length.
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}

	return string(b)
}
