package consumers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/iqbaludinm/hr-microservice/profile-service/model/domain"
	"github.com/iqbaludinm/hr-microservice/profile-service/model/kafkamodel"
	"github.com/iqbaludinm/hr-microservice/profile-service/repository"
	"go.uber.org/zap"
)

type KafkaUserConsumerService interface {
	Insert(message []byte) error
	Update(message []byte) error
	UpdatePass(message []byte) error
	// Delete(message []byte) error
}

type kafkaUserConsumerService struct {
	profileRepository repository.ProfileRepository
	logger            *zap.SugaredLogger
}

func NewKafkaUserConsumerService(profileRepository repository.ProfileRepository, logger *zap.SugaredLogger) KafkaUserConsumerService {
	return &kafkaUserConsumerService{
		profileRepository: profileRepository,
		logger:            logger,
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
	if err := s.profileRepository.CreateUser(context.TODO(), user); err != nil {
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
	if _, err := s.profileRepository.UpdateMyProfileTx(context.TODO(), user.ID, user); err != nil {
		s.logger.Errorw("error kafka update user consumer:", "error", err.Error())
	}

	return nil
}

func (s *kafkaUserConsumerService) UpdatePass(message []byte) error {
	userMsg := new(kafkamodel.KafkaUserMessage)

	if err := json.Unmarshal(message, userMsg); err != nil {
		s.logger.Errorw("error kafka update user consumer:", "error", err.Error())
	}

	user := domain.User{
		ID: userMsg.ID,
		Password:  userMsg.Password,
		UpdatedAt: time.Now(),
	}

	// update user
	if err := s.profileRepository.UpdatePasswordTx(context.TODO(), user); err != nil {
		s.logger.Errorw("error kafka update user consumer:", "error", err.Error())
	}

	return nil
}

// func (s *kafkaUserConsumerService) Delete(message []byte) error {
// 	userMsg := new(kafkamodel.KafkaUserMessage)

// 	if err := json.Unmarshal(message, userMsg); err != nil {
// 		s.logger.Errorw("error kafka delete user consumer:", "error", err.Error())
// 	}

// 	// delete user
// 	if err := s.userRepository.Delete(context.TODO(), userMsg.ID); err != nil {
// 		s.logger.Errorw("error kafka delete user consumer:", "error", err.Error())
// 	}

// 	return nil
// }
