package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/iqbaludinm/hr-microservice/auth-service/utils"
)

var logger = utils.NewLogger()

var (
	KafkaHost              = utils.GetEnv("KAFKA_HOST")
	KafkaPort              = utils.GetEnv("KAFKA_PORT")
	KafkaTopic             = utils.GetEnv("KAFKA_TOPIC")
	KafkaTopicLog          = utils.GetEnv("KAFKA_TOPIC_LOG")
	KafkaTopicNotification = utils.GetEnv("KAFKA_TOPIC_NOTIFICATION")
	KafkaSubscribeTopics   = strings.Split(utils.GetEnv("KAFKA_SUBSCRIBE_TOPICS"), ",")
	KafkaConsumerGroup     = utils.GetEnv("KAFKA_CONSUMER_GROUP")
	KafkaAddressFamily     = utils.GetEnv("KAFKA_ADDRESS_FAMILY")
	KafkaSessionTimeout, _ = strconv.Atoi(utils.GetEnv("KAFKA_SESSION_TIMEOUT"))
	KafkaAutoOffsetReset   = utils.GetEnv("KAFKA_AUTO_OFFSET_RESET")
)

// Bikin Kafka Producer
func NewKafkaProducer() *kafka.Producer {

	broker := fmt.Sprintf("%s:%s", KafkaHost, KafkaPort) // inisiasi alamat server, in case: localhost:9092

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		logger.Errorf("Error on creating kafka producer: %s\n", err.Error())
	}

	return producer
}

// Bikin Kafka Consumer
func NewKafkaConsumer() *kafka.Consumer {
	broker := fmt.Sprintf("%s:%s", KafkaHost, KafkaPort)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     broker,
		"broker.address.family": KafkaAddressFamily,
		"group.id":              KafkaConsumerGroup,
		"session.timeout.ms":    KafkaSessionTimeout,
		"auto.offset.reset":     KafkaAutoOffsetReset,
	})
	if err != nil {
		logger.Errorf("Error on creating kafka consumer: %s\n", err.Error())
	}

	logger.Info("Kafka consumer created")

	return consumer
}
