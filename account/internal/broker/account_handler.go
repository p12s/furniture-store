package broker

import (
	"encoding/json"
	"fmt"

	"github.com/p12s/furniture-store/account/internal/domain"
	"github.com/p12s/furniture-store/account/internal/service"
)

func (k *Kafka) readPayload(payload interface{}, target interface{}) error {
	jsonString, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling event value to json string: %w", err)
	}

	err = json.Unmarshal(jsonString, &target)
	if err != nil {
		return fmt.Errorf("error unmarshaling event value to []byte: %w", err)
	}

	return nil
}

func (k *Kafka) createAccount(payload interface{}, service *service.Service) error {
	var account domain.Account
	err := k.readPayload(payload, &account)
	if err != nil {
		return fmt.Errorf("read create account payload: %w/n", err)
	}

	return service.Accounter.CreateAccount(domain.Account{
		PublicId: account.PublicId,
		Name:     account.Name,
		Username: account.Username,
		Email:    account.Email,
		Address:  account.Address,
	})
}

func (k *Kafka) updateAccountInfo(payload interface{}, service *service.Service) error {
	var data domain.UpdateAccountInput
	err := k.readPayload(payload, &data)
	if err != nil {
		return fmt.Errorf("read update account info payload: %w/n", err)
	}

	return service.Accounter.UpdateAccountInfo(domain.UpdateAccountInput{
		PublicId: data.PublicId,
		Name:     data.Name,
		Username: data.Username,
		Password: data.Password,
		Email:    data.Email,
		Address:  data.Address,
	})
}

func (k *Kafka) updateAccountRole(payload interface{}, service *service.Service) error {
	var data domain.UpdateAccountRoleInput
	err := k.readPayload(payload, &data)
	if err != nil {
		return fmt.Errorf("read update account role payload: %w/n", err)
	}

	return service.Accounter.UpdateAccountRole(domain.UpdateAccountRoleInput{
		PublicId: data.PublicId,
		Role:     data.Role,
	})
}

func (k *Kafka) updateAccountToken(payload interface{}, service *service.Service) error {
	var data domain.UpdateAccountRoleInput // TODO тело не свое пока
	err := k.readPayload(payload, &data)
	if err != nil {
		return fmt.Errorf("read update account role payload: %w/n", err)
	}

	return service.Accounter.UpdateAccountRole(domain.UpdateAccountRoleInput{
		PublicId: data.PublicId,
		Role:     data.Role,
	})
}

func (k *Kafka) deleteAccount(payload interface{}, service *service.Service) error {
	var data domain.DeleteAccountInput
	err := k.readPayload(payload, &data)
	if err != nil {
		return fmt.Errorf("read account delete payload: %w/n", err)
	}

	return service.Accounter.DeleteAccount(data.PublicId)
}
