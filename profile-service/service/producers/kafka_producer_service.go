package producers

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type KafkaProducerService interface {
	Produce(data interface{}, action, topic string) error
}

type kafkaProducerService struct {
	producer *kafka.Producer
	logger   *zap.SugaredLogger
}

func NewKafkaProducerService(producer *kafka.Producer, logger *zap.SugaredLogger) KafkaProducerService {
	return &kafkaProducerService{
		producer: producer,
		logger:   logger,
	}
}

func (s *kafkaProducerService) Produce(data interface{}, action, topic string) error {
	value := new(bytes.Buffer)
	json.NewEncoder(value).Encode(data)

	s.producer.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value.Bytes(),
		Headers:        []kafka.Header{{Key: "method", Value: []byte(action)}},
	}

	for e := range s.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			m := ev
			if m.TopicPartition.Error != nil {
				s.logger.Errorw("KAFKA Delivery failed",
					"topic", *m.TopicPartition.Topic)
			} else {
				s.logger.Infow("KAFKA Delivered message",
					"topic", *m.TopicPartition.Topic,
					"partition", m.TopicPartition.Partition,
					"offset", m.TopicPartition.Offset,
				)
			}

		default:
			if strings.Contains(ev.String(), "failed:") {
				err := fiber.NewError(fiber.StatusInternalServerError, "kafka error")
				s.logger.Errorln("KAFKA Message broker error")
				return err
			} else {
				s.logger.Infow("KAFKA Ignored event",
					"event", ev)
			}
		}
	}

	return nil
}
