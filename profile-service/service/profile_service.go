package service

import (
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/iqbaludinm/hr-microservice/profile-service/config"
	"github.com/iqbaludinm/hr-microservice/profile-service/exception"
	"github.com/iqbaludinm/hr-microservice/profile-service/model/domain"
	"github.com/iqbaludinm/hr-microservice/profile-service/model/kafkamodel"
	"github.com/iqbaludinm/hr-microservice/profile-service/model/web"
	"github.com/iqbaludinm/hr-microservice/profile-service/repository"
	"github.com/iqbaludinm/hr-microservice/profile-service/service/producers"
	"go.uber.org/zap"
)

type ProfileServiceImpl struct {
	ProfileRepository repository.ProfileRepository
	Validate       *validator.Validate
}

type ProfileService interface {
	UpdateMyProfile(ctx *fiber.Ctx, id string, request web.UpdateProfileRequest) (domain.User, error)
}

type profileService struct {
	profileRepository       repository.ProfileRepository
	kafkaProducerService producers.KafkaProducerService
	logger               *zap.SugaredLogger
}

func NewProfileService(profileRepository repository.ProfileRepository, kafkaProducerService producers.KafkaProducerService, logger *zap.SugaredLogger) ProfileService {
	return &profileService{
		profileRepository:       profileRepository,
		kafkaProducerService: kafkaProducerService,
	}
}

func (service *profileService) UpdateMyProfile(ctx *fiber.Ctx, id string, request web.UpdateProfileRequest) (domain.User, error) {

	_, err := service.profileRepository.FindUserNotDeleteByQueryTx(ctx.Context(), "id", id)
	if err != nil {
		return domain.User{}, exception.ErrNotFound("User not found.")
	}

	updatedUser := domain.User{
		Name:      request.Name,
		Email:     request.Email,
		Phone:     request.Phone,
		UpdatedAt: time.Now(),
	}

	data, err := service.profileRepository.UpdateMyProfileTx(ctx.Context(), id, updatedUser)

	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			if strings.Contains(err.Error(), "users_email_key") {
				return domain.User{}, exception.ErrBadRequest("Email already exist.")
			} else if strings.Contains(err.Error(), "users_phone_key") {
				return domain.User{}, exception.ErrBadRequest("Phone already exist.")
			}
		}
		return domain.User{}, err
	}

	updatedUser.ID = id
	// produce to kafka
	kafkaRegisterMessage := kafkamodel.NewKafkaUserMessage(updatedUser)
	go service.kafkaProducerService.Produce(kafkaRegisterMessage, "PUT.USER", config.KafkaTopic)

	return data, nil
}