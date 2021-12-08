package broker

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/p12s/furniture-store/account/internal/config"
	"github.com/p12s/furniture-store/account/internal/domain"
	"github.com/sirupsen/logrus"
)

const (
	SECURITY_PROTOCOL = "SASL_SSL"
	SASL_MECHANISMS   = "SCRAM-SHA-256"
)

var _ Producer = (*BrokerProduce)(nil)

type Producer interface {
	Produce(evetType domain.EventType, eventTopic string, eventPayload interface{}) error
}

type BrokerProduce struct {
	connection *kafka.Producer
	// TopicAccountBE, TopicAccountCUD   string
	// TopicProductBE, TopicProductCUD   string
	// TopicOrderBE, TopicOrderCUD       string
	// TopicDeliveryBE, TopicDeliveryCUD string
	// TopicBillingBE, TopicBillingCUD   string
}

func NewProducer(conf *config.Cloudkarafka) (*BrokerProduce, error) { // ???? return error
	connection, err := kafka.NewProducer(&kafka.ConfigMap{
		"metadata.broker.list": conf.Brokers,
		"security.protocol":    SECURITY_PROTOCOL,
		"sasl.mechanisms":      SASL_MECHANISMS,
		"sasl.username":        conf.Username,
		"sasl.password":        conf.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("create kafka producer fail: %w", err)
	}

	return &BrokerProduce{
		connection: connection,
		// TopicAccountBE:   conf.TopicAccountBE,
		// TopicAccountCUD:  conf.TopicAccountCUD,
		// TopicProductBE:   conf.TopicProductBE,
		// TopicProductCUD:  conf.TopicProductCUD,
		// TopicOrderBE:     conf.TopicOrderBE,
		// TopicOrderCUD:    conf.TopicOrderCUD,
		// TopicDeliveryBE:  conf.TopicDeliveryBE,
		// TopicDeliveryCUD: conf.TopicDeliveryBE,
		// TopicBillingBE:   conf.TopicBillingBE,
		// TopicBillingCUD:  conf.TopicBillingCUD,
	}, nil
}

func (k *BrokerProduce) Produce(evetType domain.EventType, eventTopic string, eventPayload interface{}) error {
	deliveryChan := make(chan kafka.Event)

	var data bytes.Buffer
	if err := json.NewEncoder(&data).Encode(domain.Event{
		Type:  evetType,
		Value: eventPayload,
	}); err != nil {
		return fmt.Errorf("event encode fail: %w/n", err)
	}

	err := k.connection.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &eventTopic,
			Partition: kafka.PartitionAny,
		},
		Value: data.Bytes(),
	}, deliveryChan)
	if err != nil {
		return fmt.Errorf("event produce fail: %w/n", err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery topic-partition fail: %w/n", m.TopicPartition.Error)
	} else {
		logrus.Printf("delivered message to topic %s [%d] at offset %v/n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)

	return nil
}
