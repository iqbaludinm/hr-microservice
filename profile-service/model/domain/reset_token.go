package domain

import "time"

type ResetPasswordToken struct {
	Id          string    `json:"id"`
	Tokens      string    `json:"tokens"`
	Email       string    `json:"email"`
	URL         string    `json:"url"`
	Attempt     *int      `json:"attempt"`
	LastAttempt time.Time `json:"last_attempt"`
}
