package web

import "time"

type LogCreateRequest struct {
	Actor     string    `json:"actor"`
	ActorName string    `json:"actor_name"`
	Category  string    `json:"category"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}
