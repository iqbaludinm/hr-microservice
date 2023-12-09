package consumers

import (
	"context"
	"encoding/json"

	"github.com/iqbaludinm/hr-microservice/auth-service/model/domain"
	"github.com/iqbaludinm/hr-microservice/auth-service/model/kafkamodel"
	"github.com/iqbaludinm/hr-microservice/auth-service/repository"
	"go.uber.org/zap"
)

type KafkaAuthConsumerService interface {
	Insert(message []byte) error
	Update(message []byte) error
	// Delete(message []byte) error
}

type kafkaAuthConsumerService struct {
	authRepository repository.AuthRepository
	logger         *zap.SugaredLogger
}

func NewKafkaUserConsumerService(authRepository repository.AuthRepository, logger *zap.SugaredLogger) KafkaAuthConsumerService {
	return &kafkaAuthConsumerService{
		authRepository: authRepository,
		logger:         logger,
	}
}

func (s *kafkaAuthConsumerService) Insert(message []byte) error {
	userMsg := new(kafkamodel.KafkaUserMessage)

	if err := json.Unmarshal(message, userMsg); err != nil {
		s.logger.Errorw("error kafka insert user consumer:", "error", err.Error())
	}

	user := domain.User{
		ID:          userMsg.ID,
		Name:         userMsg.Name,
		Email:        userMsg.Email,
		Password:     userMsg.Password,
		Phone:        userMsg.Phone,
		CreatedAt:    userMsg.CreatedAt,
		UpdatedAt:    userMsg.UpdatedAt,
		DeletedAt:    userMsg.DeletedAt,
	}

	// register user
	if _, err := s.authRepository.RegisterTx(context.TODO(), user); err != nil {
		s.logger.Errorw("error kafka insert user consumer:", "error", err.Error())
	}

	return nil
}

func (s *kafkaAuthConsumerService) Update(message []byte) error {
	userMsg := new(kafkamodel.KafkaUserMessage)

	if err := json.Unmarshal(message, userMsg); err != nil {
		s.logger.Errorw("error kafka update user consumer:", "error", err.Error())
	}

	user := domain.User{
		ID:          userMsg.ID,
		Name:         userMsg.Name,
		Email:        userMsg.Email,
		Password:     userMsg.Password,
		Phone:        userMsg.Phone,
		CreatedAt:    userMsg.CreatedAt,
		UpdatedAt:    userMsg.UpdatedAt,
		DeletedAt:    userMsg.DeletedAt,
	}

	// update user
	if err := s.authRepository.UpdatePasswordTx(context.TODO(), user); err != nil {
		s.logger.Errorw("error kafka update user consumer:", "error", err.Error())
	}

	return nil
}

// func (s *kafkaAuthConsumerService) Delete(message []byte) error {
// 	userMsg := new(kafkamodel.KafkaUserMessage)

// 	if err := json.Unmarshal(message, userMsg); err != nil {
// 		s.logger.Errorw("error kafka delete user consumer:", "error", err.Error())
// 	}

// 	// delete user
// 	if err := s.authRepository.Delete(context.TODO(), userMsg.ID); err != nil {
// 		s.logger.Errorw("error kafka delete user consumer:", "error", err.Error())
// 	}

// 	return nil
// }
