package segment

import "time"

// SegmentMessage fields common to all.
type SegmentMessage struct {
	MessageId    string                 `json:"messageId"`
	Timestamp    time.Time              `json:"timestamp"`
	SentAt       time.Time              `json:"sentAt,omitempty"`
	ProjectId    string                 `json:"projectId"`
	Type         string                 `json:"type"`
	Context      map[string]interface{} `json:"context,omitempty"` // Duplicate here for batch
	Properties   map[string]interface{} `json:"properties,omitempty"`
	Traits       map[string]interface{} `json:"traits,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"` // Probably won't use
	AnonymousId  string                 `json:"anonymousId,omitempty"`
	UserId       string                 `json:"userId,omitempty"`
	Event        string                 `json:"event,omitempty"`    // Track only
	Category     string                 `json:"category,omitempty"` // Page only
	Name         string                 `json:"name,omitempty"`     // Page only
}

// SegmentBatch contains batch of messages
type SegmentBatch struct {
	MessageId string                 `json:"messageId,omitempty"`
	Timestamp time.Time              `json:"timestamp,omitempty"`
	SentAt    time.Time              `json:"sentAt,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Messages  []SegmentMessage       `json:"batch"`
}

// SegmentEvent is single message with write key
type SegmentEvent struct {
	WriteKey string `json:"writeKey,omitempty"` // Read clear, and set proejctId
	SegmentMessage
}
