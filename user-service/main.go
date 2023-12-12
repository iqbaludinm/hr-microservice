package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	// "github.com/ansrivas/fiberprometheus/v2"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/iqbaludinm/hr-microservice/user-service/config"
	"github.com/iqbaludinm/hr-microservice/user-service/controller"
	"github.com/iqbaludinm/hr-microservice/user-service/exception"
	"github.com/iqbaludinm/hr-microservice/user-service/repository"
	"github.com/iqbaludinm/hr-microservice/user-service/repository/query"
	"github.com/iqbaludinm/hr-microservice/user-service/service"
	"github.com/iqbaludinm/hr-microservice/user-service/service/consumers"
	"github.com/iqbaludinm/hr-microservice/user-service/service/producers"
	"github.com/iqbaludinm/hr-microservice/user-service/utils"
	"go.uber.org/zap"
)

// go:embed templates
var templateFS embed.FS

var logger = utils.NewLogger()

func controllers() {
	time.Local = time.UTC

	app := fiber.New(fiber.Config{
		ErrorHandler:          exception.ErrorHandler,
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	app.Use(requestid.New()) // ini middleware bawaan dari go-fiber untuk identifikasi
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
		Fields: []string{"requestId", "ip", "status", "method", "latency", "path", "resBody", "ua"},
	}))

	serverConfig := config.NewServerConfig()
	db := config.NewPostgresDatabase()
	store := repository.NewStore(db)
	validate := validator.New()

	kafkaProducer := config.NewKafkaProducer()
	kafkaProducerService := producers.NewKafkaProducerService(kafkaProducer, logger.Sugar())

	userQuery := query.NewUser()
	userRepository := repository.NewUser(store, userQuery)
	userService := service.NewUserService(userRepository, kafkaProducerService, logger.Sugar())
	userController := controller.NewUserController(validate, kafkaProducerService, userService)

	userController.Route(app)

	err := app.Listen(serverConfig.Host)
	if err != nil {
		sugar.Fatal(err)
	}
	log.Println(err)

}

func main() {
	time.Local = time.UTC

	go controllers()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	consumer := config.NewKafkaConsumer()

	if err := consumer.SubscribeTopics(config.KafkaSubscribeTopics, nil); err != nil {
		logger.Errorw("Error on subscribe topics", "error", err.Error())
	} else {
		logger.Infoln("Kafka subscribe to topics:", strings.Join(config.KafkaSubscribeTopics, ", "))
	}

	db := config.NewPostgresDatabase()
	store := repository.NewStore(db)
	userQuery := query.NewUser()
	userRepository := repository.NewUser(store, userQuery)
	kafkaUserConsumerService := consumers.NewKafkaUserConsumerService(userRepository, logger)

	run := true

	for run {
		select {
		case sig := <-sigchan:
			logger.Panicw("Caught signal, terminating", "signal", sig)
			run = false
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				logger.Infow("New Kafka Message",
					"Message Topic", e.TopicPartition,
					"Headers", e.Headers,
				)

				method := fmt.Sprintf("%v", e.Headers)

				switch method {

				// USER
				case `[method="POST.USER"]`:
					err := kafkaUserConsumerService.Insert(e.Value)
					if err != nil {
						logger.Panic(err)
					}
				case `[method="PUT.USER"]`:
					err := kafkaUserConsumerService.Update(e.Value)
					if err != nil {
						logger.Panic(err)
					}
				case `[method="PUT.USER_PASS"]`:
					err := kafkaUserConsumerService.UpdatePass(e.Value)
					if err != nil {
						logger.Panic(err)
					}
				case `[method="DELETE.USER"]`:
					err := kafkaUserConsumerService.Delete(e.Value)
					if err != nil {
						logger.Panic(err)
					}
				}

			case kafka.Error:
				logger.Errorw("Kafka Error",
					"Error", e,
					"Error_Code", e.Code(),
				)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
					logger.Panic(e.Code())
				}
			default:
				logger.Infow("Kafka Info: Ignored Event", "Event", e)
			}
		}
	}

	consumer.Close()
	logger.Info("Kafka consumer closed")
}
