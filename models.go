package main

import "encoding/json"

type EventRequest struct {
	EventType string          `json:"event_type"`
	Source    string          `json:"source"`
	Payload   json.RawMessage `json:"payload"`
}
