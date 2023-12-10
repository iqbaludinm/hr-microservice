package consumers

import (
	"context"
	"encoding/json"

	"github.com/iqbaludinm/hr-microservice/user-service/model/domain"
	"github.com/iqbaludinm/hr-microservice/user-service/model/kafkamodel"
	"github.com/iqbaludinm/hr-microservice/user-service/repository"
	"go.uber.org/zap"
)

type KafkaUserConsumerService interface {
	Insert(message []byte) error
	Update(message []byte) error
	Delete(message []byte) error
}

type kafkaUserConsumerService struct {
	userRepository repository.UserRepository
	logger         *zap.SugaredLogger
}

func NewKafkaUserConsumerService(userRepository repository.UserRepository, logger *zap.SugaredLogger) KafkaUserConsumerService {
	return &kafkaUserConsumerService{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (s *kafkaUserConsumerService) Insert(message []byte) error {
	userMsg := new(kafkamodel.KafkaUserMessage)

	if err := json.Unmarshal(message, userMsg); err != nil {
		s.logger.Errorw("error kafka insert user consumer:", "error", err.Error())
	}

	user := domain.User{
		ID:        userMsg.ID,
		Name:      userMsg.Name,
		Email:     userMsg.Email,
		Password:  userMsg.Password,
		Phone:     userMsg.Phone,
		CreatedAt: userMsg.CreatedAt,
		UpdatedAt: userMsg.UpdatedAt,
		DeletedAt: userMsg.DeletedAt,
	}

	// create user
	if err := s.userRepository.CreateUser(context.TODO(), user); err != nil {
		s.logger.Errorw("error kafka insert user consumer:", "error", err.Error())
	}

	return nil
}

func (s *kafkaUserConsumerService) Update(message []byte) error {
	userMsg := new(kafkamodel.KafkaUserMessage)

	if err := json.Unmarshal(message, userMsg); err != nil {
		s.logger.Errorw("error kafka update user consumer:", "error", err.Error())
	}

	user := domain.User{
		ID:        userMsg.ID,
		Name:      userMsg.Name,
		Email:     userMsg.Email,
		Phone:     userMsg.Phone,
		UpdatedAt: userMsg.UpdatedAt,
	}

	// update user
	if err := s.userRepository.UpdateUser(context.TODO(), user.ID, user); err != nil {
		s.logger.Errorw("error kafka update user consumer:", "error", err.Error())
	}

	return nil
}

func (s *kafkaUserConsumerService) Delete(message []byte) error {
	userMsg := new(kafkamodel.KafkaUserMessage)

	if err := json.Unmarshal(message, userMsg); err != nil {
		s.logger.Errorw("error kafka delete user consumer:", "error", err.Error())
	}

	// delete user
	if err := s.userRepository.Delete(context.TODO(), userMsg.ID); err != nil {
		s.logger.Errorw("error kafka delete user consumer:", "error", err.Error())
	}

	return nil
}
