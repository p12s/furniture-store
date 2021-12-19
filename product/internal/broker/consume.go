package broker

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/p12s/furniture-store/product/internal/config"
	"github.com/p12s/furniture-store/product/internal/domain"
	"github.com/p12s/furniture-store/product/internal/service"
	"github.com/sirupsen/logrus"
)

const (
	AUTO_OFFSET_RESET = "earliest"
)

var _ Consumer = (*BrokerConsume)(nil)

type Consumer interface {
	Subscribe() error
	ProcessEvent(event domain.Event)
}

type BrokerConsume struct {
	connection                        *kafka.Consumer
	service                           *service.Service
	TopicAccountBE, TopicAccountCUD   string
	TopicProductBE, TopicProductCUD   string
	TopicOrderBE, TopicOrderCUD       string
	TopicDeliveryBE, TopicDeliveryCUD string
	TopicBillingBE, TopicBillingCUD   string
}

func NewConsumer(service *service.Service, conf *config.Broker) (*BrokerConsume, error) {
	connection, err := kafka.NewConsumer(&kafka.ConfigMap{
		"metadata.broker.list": conf.Brokers,
		"security.protocol":    SECURITY_PROTOCOL,
		"sasl.mechanisms":      SASL_MECHANISMS,
		"sasl.username":        conf.Username,
		"sasl.password":        conf.Password,
		"group.id":             conf.GroupId,
		"auto.offset.reset":    AUTO_OFFSET_RESET,
	})
	if err != nil {
		return nil, fmt.Errorf("create kafka consumer fail: %w", err)
	}

	return &BrokerConsume{
		connection:       connection,
		service:          service,
		TopicAccountBE:   conf.TopicAccountBE,
		TopicAccountCUD:  conf.TopicAccountCUD,
		TopicProductBE:   conf.TopicProductBE,
		TopicProductCUD:  conf.TopicProductCUD,
		TopicOrderBE:     conf.TopicOrderBE,
		TopicOrderCUD:    conf.TopicOrderCUD,
		TopicDeliveryBE:  conf.TopicDeliveryBE,
		TopicDeliveryCUD: conf.TopicDeliveryCUD,
		TopicBillingBE:   conf.TopicBillingBE,
		TopicBillingCUD:  conf.TopicBillingCUD,
	}, nil
}

func (k *BrokerConsume) Subscribe() error {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	topics := []string{
		k.TopicAccountBE, k.TopicAccountCUD,
		k.TopicProductBE, k.TopicProductCUD,
		k.TopicOrderBE, k.TopicOrderCUD,
		k.TopicDeliveryBE, k.TopicDeliveryCUD,
		k.TopicBillingBE, k.TopicBillingCUD,
	}
	fmt.Printf("topics: %T %v /n", topics, topics)

	err := k.connection.SubscribeTopics([]string{
		k.TopicAccountBE, k.TopicAccountCUD,
		k.TopicProductBE, k.TopicProductCUD,
		k.TopicOrderBE, k.TopicOrderCUD,
		k.TopicDeliveryBE, k.TopicDeliveryCUD,
		k.TopicBillingBE, k.TopicBillingCUD,
	}, nil)
	if err != nil {
		return fmt.Errorf("subscribe broker topics fail: %w", err)
	}

	run := true
	for run == true { // nolint
		select {
		case sig := <-sigchan:
			logrus.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := k.connection.ReadMessage(1 * time.Second)
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
			k.ProcessEvent(eventData)
		}
	}

	logrus.Println("closing consumer")
	err = k.connection.Close()
	if err != nil {
		return fmt.Errorf("closing consumer fail: %w", err)
	}
	return nil
}

func (k *BrokerConsume) ProcessEvent(event domain.Event) {
	switch event.Type {
	case domain.EVENT_ACCOUNT_CREATED:
		err := k.createAccount(event.Value)
		if err != nil {
			logrus.Errorf("process 'create account' event fail: %s/n", err.Error())
		}
	case domain.EVENT_ACCOUNT_INFO_UPDATED:
		err := k.updateAccountInfo(event.Value)
		if err != nil {
			logrus.Errorf("process 'update account info' event fail: %s/n", err.Error())
		}
	case domain.EVENT_ACCOUNT_ROLE_UPDATED:
		err := k.updateAccountRole(event.Value)
		if err != nil {
			logrus.Errorf("process 'update account role' event fail: %s/n", err.Error())
		}
	case domain.EVENT_ACCOUNT_TOKEN_UPDATED:
		err := k.updateAccountToken(event.Value)
		if err != nil {
			logrus.Errorf("process 'update account token' event fail: %s/n", err.Error())
		}
	case domain.EVENT_ACCOUNT_DELETED:
		err := k.deleteAccount(event.Value)
		if err != nil {
			logrus.Errorf("process 'delete account' event fail: %s/n", err.Error())
		}
	default:
		fmt.Printf("unknown event type: %v/n", event.Value)
	}
}

func (k *BrokerConsume) createAccount(payload interface{}) error {
	var account domain.Account
	err := readPayload(payload, &account)
	if err != nil {
		return fmt.Errorf("account-create payload fail: %w/n", err)
	}

	return k.service.CreateAccount(domain.Account{
		PublicId: account.PublicId,
		Name:     account.Name,
		Username: account.Username,
		Email:    account.Email,
		Address:  account.Address,
	})
}

func (k *BrokerConsume) updateAccountInfo(payload interface{}) error {
	var data domain.UpdateAccountInput
	err := readPayload(payload, &data)
	if err != nil {
		return fmt.Errorf("account-update info payload fail: %w/n", err)
	}

	return k.service.UpdateAccountInfo(domain.UpdateAccountInput{
		PublicId: data.PublicId,
		Name:     data.Name,
		Username: data.Username,
		Password: data.Password,
		Email:    data.Email,
		Address:  data.Address,
	})
}

func (k *BrokerConsume) updateAccountRole(payload interface{}) error {
	var data domain.UpdateAccountRoleInput
	err := readPayload(payload, &data)
	if err != nil {
		return fmt.Errorf("account-role update payload fail: %w/n", err)
	}

	return k.service.UpdateAccountRole(domain.UpdateAccountRoleInput{
		PublicId: data.PublicId,
		Role:     data.Role,
	})
}

func (k *BrokerConsume) updateAccountToken(payload interface{}) error {
	var data domain.UpdateAccountRoleInput // TODO тело не свое пока
	err := readPayload(payload, &data)
	if err != nil {
		return fmt.Errorf("account-role update payload fail: %w/n", err)
	}

	return k.service.UpdateAccountRole(domain.UpdateAccountRoleInput{
		PublicId: data.PublicId,
		Role:     data.Role,
	})
}

func (k *BrokerConsume) deleteAccount(payload interface{}) error {
	var data domain.DeleteAccountInput
	err := readPayload(payload, &data)
	if err != nil {
		return fmt.Errorf("delete-account payload fail: %w/n", err)
	}

	return k.service.DeleteAccount(data.PublicId)
}

func readPayload(payload interface{}, target interface{}) error {
	jsonString, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling event value to json string fail: %w", err)
	}

	err = json.Unmarshal(jsonString, &target)
	if err != nil {
		return fmt.Errorf("unmarshaling event value to []byte fail: %w", err)
	}

	return nil
}
