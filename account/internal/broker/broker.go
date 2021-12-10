package broker

import (
	"fmt"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/p12s/furniture-store/account/internal/config"
	"github.com/p12s/furniture-store/account/internal/service"
)

//go:generate mockgen -destination mocks/mock.go -package broker github.com/p12s/furniture-store/account/internal/broker Consumer,Producer

// Broker
type Broker struct {
	Producer
	Consumer
	TopicAccountBE, TopicAccountCUD   string
	TopicProductBE, TopicProductCUD   string
	TopicOrderBE, TopicOrderCUD       string
	TopicDeliveryBE, TopicDeliveryCUD string
	TopicBillingBE, TopicBillingCUD   string
}

// NewBroker - constructor
func NewBroker(service *service.Service, config *config.Broker) (*Broker, error) {
	producer, err := NewProducer(config)
	if err != nil {
		return nil, fmt.Errorf("broker producer fail: %w/n", err)
	}
	consumer, err := NewConsumer(service, config)
	if err != nil {
		return nil, fmt.Errorf("broker consumer fail: %w/n", err)
	}

	return &Broker{
		Producer:         producer,
		Consumer:         consumer,
		TopicAccountBE:   config.TopicAccountBE,
		TopicAccountCUD:  config.TopicAccountCUD,
		TopicProductBE:   config.TopicProductBE,
		TopicProductCUD:  config.TopicProductCUD,
		TopicOrderBE:     config.TopicOrderBE,
		TopicOrderCUD:    config.TopicOrderCUD,
		TopicDeliveryBE:  config.TopicDeliveryBE,
		TopicDeliveryCUD: config.TopicDeliveryCUD,
		TopicBillingBE:   config.TopicBillingBE,
		TopicBillingCUD:  config.TopicBillingCUD,
	}, nil
}
