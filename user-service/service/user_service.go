package service

import (
	"context"
	"log"
	"strings"

	// "fmt"
	// "strings"
	// "time"

	// "github.com/google/uuid"
	// "github.com/iqbaludinm/hr-microservice/user-service/config"
	"github.com/iqbaludinm/hr-microservice/user-service/exception"
	"github.com/iqbaludinm/hr-microservice/user-service/model/domain"

	// "github.com/iqbaludinm/hr-microservice/user-service/model/kafkamodel"
	"github.com/iqbaludinm/hr-microservice/user-service/model/web"
	"github.com/iqbaludinm/hr-microservice/user-service/repository"
	"github.com/iqbaludinm/hr-microservice/user-service/service/producers"
	"go.uber.org/zap"
)

type UserService interface {
	// With Transaction
	// CreateUser(ctx context.Context, request web.CreateUserRequest) (web.UserResponse, error)
	// UpdateUser(ctx context.Context, id string, request web.UpdateUserRequest) (web.UserResponse, error)
	// UpdatePassword(ctx context.Context, request web.UpdateUserPasswordRequest) (web.UserResponse, error)
	// Delete(ctx context.Context, id string) (web.UserResponse, error)

	// Without Transaction
	FindAllUser(ctx context.Context, filter web.UserQueryFilter) (result []web.UserResponse, totalData int, err error)
	// FindById(ctx context.Context, id string, filter web.UserQueryFilter) (web.UserResponse, error)
	// FindByEmail(ctx context.Context, email string) (web.UserResponse, error)
	// FindByPhoneNumber(ctx context.Context, phone string) (web.UserResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
	kafkaProducerService producers.KafkaProducerService
	logger *zap.SugaredLogger
}

func NewUserService(userRepository repository.UserRepository, kafkaProducerService producers.KafkaProducerService, logger *zap.SugaredLogger) UserService {
	return &userService{
		userRepository: userRepository,
		kafkaProducerService: kafkaProducerService,
	}
}

// func (s *userService) CreateUser(c context.Context, request web.CreateUserRequest) (web.UserResponse, error) {
// 	// validate email exist
// 	if _, err := s.userRepository.FindByEmail(c, request.Email); err != nil {
// 		if strings.Contains(err.Error(), "no rows") {
// 			return web.UserResponse{}, exception.ErrNotFound(fmt.Sprintf("Supervisor %s not found", request.Email))
// 		} else {
// 			return web.UserResponse{}, err
// 		}
// 	}

// 	// convert to domain or model user
// 	user := domain.User{
// 		ID: uuid.New().String(),
// 		Name: request.Name,
// 		Email: request.Email,
// 		Password: request.Password,
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}

// 	// call the repo for inserting to db
// 	if err := s.userRepository.CreateUser(c, user); err != nil {
// 		s.logger.Infow(err.Error(), "Create User Error")
// 		return web.UserResponse{}, err
// 	}

// 	// get or returning the repo have created to db
// 	newUser, err := s.userRepository.FindById(c, user.ID, domain.UserQueryFilter{
// 		ShowDeleted: false,
// 	})
// 	if err != nil {
// 		return web.UserResponse{}, exception.ErrInternalServer(fmt.Sprintf("Successfully created user, but failed to get the user have created. Error: %s", err.Error()))
// 	}

// 	// produce kafka create-user message
// 	kafkaUserMessage := kafkamodel.NewKafkaUserMessage(newUser)
// 	go s.kafkaProducerService.Produce(kafkaUserMessage, "POST.JOB", config.KafkaTopic)


// 	return newUser.
// }

// func (u *userService) UpdateUser(c context.Context, request web.CreateUserRequest) (web.UserResponse, error) {
// 	return 
// }
// func (u *userService) UpdatePassword(c context.Context, request web.CreateUserRequest) (web.UserResponse, error) {
// 	return 
// }
// func (u *userService) Delete(c context.Context, request web.CreateUserRequest) (web.UserResponse, error) {
// 	return 
// }
func (u *userService) FindAllUser(c context.Context, filter web.UserQueryFilter) (result []web.UserResponse, totalData int, err error) {
	repositoryResponse, err := u.userRepository.FindAllUser(c, domain.ToDomainUserQueryFilter(filter))
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, 0, exception.ErrNotFound("User not found")
		} else {
			return nil, 0, err
		}
	}

	// get total-data
	totalData, err = u.userRepository.CountAllUser(c, domain.ToDomainUserQueryFilter(filter))
	if err != nil {
		return nil, 0, err
	}

	// convert to web.JobResponse
	for _, user := range repositoryResponse {
		log.Println(user.Name)
		result = append(result, user.ToUserResponse())
	}
	log.Println(result)

	return result, totalData, err
}
// func (u *userService) FindById(c context.Context, request web.CreateUserRequest) (web.UserResponse, error) {
// 	return 
// }
// func (u *userService) FindByEmail(c context.Context, request web.CreateUserRequest) (web.UserResponse, error) {
// 	return 
// }
// func (u *userService) FindByPhoneNumber(c context.Context, request web.CreateUserRequest) (web.UserResponse, error) {
// 	return 
// }