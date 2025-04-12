package configs

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewProducer() (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":                     os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"acks":                                  "all",
		"enable.idempotence":                    1,
		"max.in.flight.requests.per.connection": 1,
		"retries":                               5,
		"retry.backoff.ms":                      100,
		"delivery.timeout.ms":                   1000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new producer: %v", err)
	}

	return p, nil
}

func NewConsumer() (*kafka.Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":         os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":                  "flight_group",
		"auto.offset.reset":         "earliest",
		"enable.auto.commit":        false,
		"session.timeout.ms":        6000,
		"max.poll.interval.ms":      300000,
		"heartbeat.interval.ms":     2000,
		"fetch.min.bytes":           1,
		"max.partition.fetch.bytes": 1048576,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new consumer: %v", err)
	}
	return c, err
}
