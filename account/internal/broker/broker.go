package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/p12s/furniture-store/account/internal/domain"
	"github.com/p12s/furniture-store/account/internal/service"
	"github.com/sirupsen/logrus"
)

type Kafka struct {
	Consumer *kafka.Consumer
	Producer *kafka.Producer
	Topics   map[string]string

	// TopicAccountBE   string
	// TopicAccountCUD  string
	// TopicProductBE   string
	// TopicProductCUD  string
	// TopicOrderBE     string
	// TopicOrderCUD    string
	// TopicDeliveryBE  string
	// TopicDeliveryCUD string
}

type KafkaConfig struct {
	Brokers, Username, Password, GroupId string
	Topics                               map[string]string
}

func NewKafka(conf KafkaConfig) (*Kafka, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"metadata.broker.list": conf.Brokers, // os.Getenv("CLOUDKARAFKA_BROKERS"),
		"security.protocol":    "SASL_SSL",
		"sasl.mechanisms":      "SCRAM-SHA-256",
		"sasl.username":        conf.Username, // os.Getenv("CLOUDKARAFKA_USERNAME"),
		"sasl.password":        conf.Password, // os.Getenv("CLOUDKARAFKA_PASSWORD"),
	})
	if err != nil {
		return nil, fmt.Errorf("error in kafka constructor, while create producer: %w", err)
	}

	// TODO выкинуть консьюмер - в аккаунте он не нужен
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"metadata.broker.list": conf.Brokers, // os.Getenv("CLOUDKARAFKA_BROKERS"),
		"security.protocol":    "SASL_SSL",
		"sasl.mechanisms":      "SCRAM-SHA-256",
		"sasl.username":        conf.Username, // os.Getenv("CLOUDKARAFKA_USERNAME"),
		"sasl.password":        conf.Password, // os.Getenv("CLOUDKARAFKA_PASSWORD"),
		"group.id":             conf.GroupId,  // os.Getenv("CLOUDKARAFKA_GROUP_ID"),
		"auto.offset.reset":    "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("error in kafka constructor, while create consumer: %w", err)
	}

	return &Kafka{
		Producer: producer,
		Consumer: consumer,
		Topics:   conf.Topics,
		// TopicAccountBE:  os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX") + "account",
		// TopicAccountCUD: os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX") + "stream",
		// TopicTaskBE:     os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX") + "task",
		// TopicTaskCUD:    os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX") + "stream",
		// TopicBillingBE:  os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX") + "billing",
		// TopicBillingCUD: os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX") + "stream",
	}, nil
}

// TODO зарефакторить eventTopic
func (k *Kafka) Event(evetType domain.EventType, eventTopic string, eventPayload interface{}) {
	deliveryChan := make(chan kafka.Event)

	var data bytes.Buffer
	if err := json.NewEncoder(&data).Encode(domain.Event{
		Type:  evetType,
		Value: eventPayload,
	}); err != nil {
		fmt.Printf("auth brocker data encode: %s\n", err.Error()) // TODO logrus
		return
	}

	err := k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &eventTopic,
			Partition: kafka.PartitionAny,
		},
		Value: data.Bytes(),
	}, deliveryChan)
	if err != nil {
		fmt.Printf("auth broker produce: %s\n", err.Error())
		return
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
}

func (k *Kafka) Subscribe(service *service.Service) error {
	// if err := godotenv.Load(); err != nil {
	// 	fmt.Printf("error loading env variables: %s\n", err.Error())
	// 	return fmt.Errorf("error in kafka constructor, while create consumer: %w", err)
	// }

	// topics := []string{
	// 	k.TopicAccountBE, k.TopicAccountCUD,
	// 	k.TopicTaskBE, k.TopicBillingBE,
	// }

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	err := k.Consumer.SubscribeTopics([]string{"some"}, nil) // k.Topics to array
	if err != nil {
		return fmt.Errorf("subscribe topics kafka: %w", err)
	}

	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			logrus.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := k.Consumer.ReadMessage(1 * time.Second)
			if err != nil {
				continue
			}
			fmt.Printf("✅ Message on %s:\nvalue: %s\n", ev.TopicPartition, string(ev.Value)) // TODO удалить вывод после реализации/обкатки всех событий
			var eventData domain.Event
			err = json.Unmarshal(ev.Value, &eventData)
			if err != nil {
				logrus.Errorf("Unmarshal error: %s\n", err.Error())
				continue
			}
			k.processEvent(eventData, service)
		}
	}

	logrus.Println("Closing consumer")
	k.Consumer.Close()
	return nil
}

func (k *Kafka) processEvent(event domain.Event, service *service.Service) {
	switch event.Type {
	case domain.EVENT_ACCOUNT_CREATED:
		k.createAccount(event.Value, service)
	case domain.EVENT_ACCOUNT_INFO_UPDATED:
		k.updateAccountInfo(event.Value, service)
	case domain.EVENT_ACCOUNT_ROLE_UPDATED:
		k.updateAccountRole(event.Value, service)
	case domain.EVENT_ACCOUNT_TOKEN_UPDATED:
		k.updateAccountToken(event.Value, service)
	case domain.EVENT_ACCOUNT_DELETED:
		k.deleteAccount(event.Value, service)

	default:
		fmt.Println("unknown event type")
	}
}
