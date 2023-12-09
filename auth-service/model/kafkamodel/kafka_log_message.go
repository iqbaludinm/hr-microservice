package kafkamodel

import "time"

// This struct is used for mapping the log message, and produce the data to 'kafka' 'log' topic.
type KafkaLogMessage struct {
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}
