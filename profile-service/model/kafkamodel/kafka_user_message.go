package kafkamodel

import (
	"time"

	"github.com/iqbaludinm/hr-microservice/profile-service/model/domain"
)

// This struct is used for mapping the incoming message from 'kafka' (json) to an object.
// Since we will only consume 'user' data from other service, this struct data will only be used
// to consume 'user' data from 'kafka' to the database.
// We won't use it to produce 'user' data to 'kafka'
type KafkaUserMessage struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Phone     string     `json:"phone"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// Convert "User" object to "KafkaUserMessage" object
func NewKafkaUserMessage(user domain.User) KafkaUserMessage {
	return KafkaUserMessage{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Password: user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}
