package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/p12s/furniture-store/account/internal/broker"
	mock_broker "github.com/p12s/furniture-store/account/internal/broker/mocks"
	"github.com/p12s/furniture-store/account/internal/config"
	"github.com/p12s/furniture-store/account/internal/domain"
	"github.com/p12s/furniture-store/account/internal/service"
	mock_service "github.com/p12s/furniture-store/account/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func init() {
	currentDir, err := os.Getwd()
	assert.ObjectsAreEqual(nil, err)

	configPath := filepath.Dir(filepath.Dir(filepath.Dir(currentDir)))
	err = godotenv.Load(os.ExpandEnv(fmt.Sprintf("%s/%s", configPath, DIR_ENV_PATH)))
	assert.ObjectsAreEqual(nil, err)

	_, err = config.New()
	assert.ObjectsAreEqual(nil, err)

	assert.ObjectsAreEqual("dev", os.Getenv("ENV_CURRENT"))
}

func TestHandler_updateAccount(t *testing.T) {

	type accountMockBehavior func(s *mock_service.MockAccounter, account domain.UpdateAccountInput)
	type brokerMockProducer func(s *mock_broker.MockProducer, event domain.EventType, topic string, input interface{})

	publicId, _ := uuid.Parse("265cee57-2ff9-4ed3-85e1-d3373fa2a1a5")
	var name = "Ivan"
	var username = "ivan"
	var password = "qwerty"
	var email = "test@test.ru"
	var address = "Some-city, some-street, some-hause"

	tests := []struct {
		name                string
		inputBody           string
		inputAccount        domain.UpdateAccountInput
		accountMockBehavior accountMockBehavior
		eventType           domain.EventType
		eventData           interface{}
		topic               string
		brokerMockProducer  brokerMockProducer
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "Can update account with correct input",
			inputBody: `{"public_id": "265cee57-2ff9-4ed3-85e1-d3373fa2a1a5", "name": "Ivan", "username": "ivan", "password": "qwerty", "email": "test@test.ru", "address": "Some-city, some-street, some-hause"}`,
			inputAccount: domain.UpdateAccountInput{
				PublicId: publicId,
				Name:     &name,
				Username: &username,
				Password: &password,
				Email:    &email,
				Address:  &address,
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, account domain.UpdateAccountInput) {
				s.EXPECT().UpdateAccountInfo(account).Return(nil)
			},
			eventType: domain.EVENT_ACCOUNT_INFO_UPDATED,
			eventData: domain.UpdateAccountInput{
				PublicId: publicId,
				Name:     &name,
				Username: &username,
				Password: &password,
				Email:    &email,
				Address:  &address,
			},
			topic: "",
			brokerMockProducer: func(s *mock_broker.MockProducer, event domain.EventType, topic string, input interface{}) {
				s.EXPECT().Produce(event, topic, input).Return(nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ``,
		},
		{
			name:      "Can't update account without required public account id",
			inputBody: `{"name": "Ivan", "username": "ivan", "password": "qwerty", "email": "test@test.ru", "address": "Some-city, some-street, some-hause"}`,
			inputAccount: domain.UpdateAccountInput{
				Name:     &name,
				Username: &username,
				Password: &password,
				Email:    &email,
				Address:  &address,
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, account domain.UpdateAccountInput) {},
			eventType:           domain.EVENT_ACCOUNT_INFO_UPDATED,
			eventData:           domain.UpdateAccountInput{},
			topic:               "",
			brokerMockProducer:  func(s *mock_broker.MockProducer, event domain.EventType, topic string, input interface{}) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Can return error response if service failure",
			inputBody: `{"public_id": "265cee57-2ff9-4ed3-85e1-d3373fa2a1a5", "name": "Ivan", "username": "ivan", "password": "qwerty", "email": "test@test.ru", "address": "Some-city, some-street, some-hause"}`,
			inputAccount: domain.UpdateAccountInput{
				PublicId: publicId,
				Name:     &name,
				Username: &username,
				Password: &password,
				Email:    &email,
				Address:  &address,
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, account domain.UpdateAccountInput) {
				s.EXPECT().UpdateAccountInfo(account).Return(errors.New(""))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			acc := mock_service.NewMockAccounter(ctrl)
			tt.accountMockBehavior(acc, tt.inputAccount)
			serviceMock := &service.Service{Accounter: acc}

			brokerProducer := mock_broker.NewMockProducer(ctrl)
			var brokerMock *broker.Broker
			if tt.brokerMockProducer != nil {
				tt.brokerMockProducer(brokerProducer, tt.eventType, tt.topic, tt.eventData)
				brokerMock = &broker.Broker{Producer: brokerProducer}
			}

			handler := NewHandler(serviceMock, brokerMock)
			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.PUT("/account/info", handler.updateAccount)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/account/info", bytes.NewBufferString(tt.inputBody))

			r.ServeHTTP(w, req)
			time.Sleep(WAITING_GORUTINE_END_TIME)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_deleteAccount(t *testing.T) {

	type accountMockBehavior func(s *mock_service.MockAccounter, account domain.DeleteAccountInput)
	type brokerMockProducer func(s *mock_broker.MockProducer, event domain.EventType, topic string, input interface{})

	publicId := "265cee57-2ff9-4ed3-85e1-d3373fa2a1a5"

	tests := []struct {
		name                string
		inputBody           string
		inputAccount        domain.DeleteAccountInput
		accountMockBehavior accountMockBehavior
		eventType           domain.EventType
		eventData           interface{}
		topic               string
		brokerMockProducer  brokerMockProducer
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "Can delete account with correct input",
			inputBody: `{"public_id": "265cee57-2ff9-4ed3-85e1-d3373fa2a1a5"}`,
			inputAccount: domain.DeleteAccountInput{
				PublicId: publicId,
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, account domain.DeleteAccountInput) {
				s.EXPECT().DeleteAccount(account.PublicId).Return(nil)
			},
			eventType: domain.EVENT_ACCOUNT_DELETED,
			eventData: domain.DeleteAccountInput{
				PublicId: publicId,
			},
			topic: "",
			brokerMockProducer: func(s *mock_broker.MockProducer, event domain.EventType, topic string, input interface{}) {
				s.EXPECT().Produce(event, topic, input).Return(nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ``,
		},
		{
			name:                "Can't delete account without required public account id",
			inputBody:           `{}`,
			inputAccount:        domain.DeleteAccountInput{},
			accountMockBehavior: func(s *mock_service.MockAccounter, account domain.DeleteAccountInput) {},
			eventType:           domain.EVENT_ACCOUNT_DELETED,
			eventData:           domain.DeleteAccountInput{},
			topic:               "",
			brokerMockProducer:  func(s *mock_broker.MockProducer, event domain.EventType, topic string, input interface{}) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Can return error response if service failure",
			inputBody: `{"public_id": "265cee57-2ff9-4ed3-85e1-d3373fa2a1a5"}`,
			inputAccount: domain.DeleteAccountInput{
				PublicId: publicId,
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, account domain.DeleteAccountInput) {
				s.EXPECT().DeleteAccount(account.PublicId).Return(errors.New(""))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			acc := mock_service.NewMockAccounter(ctrl)
			tt.accountMockBehavior(acc, tt.inputAccount)
			serviceMock := &service.Service{Accounter: acc}

			brokerProducer := mock_broker.NewMockProducer(ctrl)
			var brokerMock *broker.Broker
			if tt.brokerMockProducer != nil {
				tt.brokerMockProducer(brokerProducer, tt.eventType, tt.topic, tt.eventData)
				brokerMock = &broker.Broker{Producer: brokerProducer}
			}

			handler := NewHandler(serviceMock, brokerMock)
			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.DELETE("/account/", handler.deleteAccount)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/account/", bytes.NewBufferString(tt.inputBody))

			r.ServeHTTP(w, req)
			time.Sleep(WAITING_GORUTINE_END_TIME)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getAccountInfo(t *testing.T) {

	type accountMockBehavior func(s *mock_service.MockAccounter, publicId string, account domain.Account)

	publicId := "265cee57-2ff9-4ed3-85e1-d3373fa2a1a5"
	uuidPublicId, _ := uuid.Parse(publicId)

	tests := []struct {
		name                string
		publicId            string
		outputBody          string
		outputAccount       domain.Account
		accountMockBehavior accountMockBehavior
		expectedStatusCode  int
		expectedRequestBody string
		isTokenExists       bool
	}{
		{
			name:       "Can return account",
			publicId:   publicId,
			outputBody: `{"public_id": "265cee57-2ff9-4ed3-85e1-d3373fa2a1a5"}`,
			outputAccount: domain.Account{
				Id:       1,
				PublicId: uuidPublicId,
				Name:     "Ivan",
				Username: "ivan",
				Password: "qwerty",
				Email:    "test@test.ru",
				Address:  "address",
				Role:     domain.ROLE_CUSTOMER,
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, publicId string, account domain.Account) {
				s.EXPECT().GetAccount(publicId).Return(account, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"id":1,"public_id":"265cee57-2ff9-4ed3-85e1-d3373fa2a1a5","name":"Ivan","username":"ivan","password":"qwerty","email":"test@test.ru","address":"address","role":0}`,
			isTokenExists:       true,
		},
		{
			name:                "Can't return account without jwt-token in header",
			accountMockBehavior: func(s *mock_service.MockAccounter, publicId string, account domain.Account) {},
			expectedStatusCode:  http.StatusNotFound,
			expectedRequestBody: `{"message":"account public id not found"}`,
			isTokenExists:       false,
		},
		{
			name:       "Can return error response if service failure",
			publicId:   publicId,
			outputBody: `{"public_id": "265cee57-2ff9-4ed3-85e1-d3373fa2a1a5"}`,
			outputAccount: domain.Account{
				Id:       1,
				PublicId: uuidPublicId,
				Name:     "Ivan",
				Username: "ivan",
				Password: "qwerty",
				Email:    "test@test.ru",
				Address:  "address",
				Role:     domain.ROLE_CUSTOMER,
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, publicId string, account domain.Account) {
				s.EXPECT().GetAccount(publicId).Return(account, errors.New(""))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"service failure"}`,
			isTokenExists:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			acc := mock_service.NewMockAccounter(ctrl)
			tt.accountMockBehavior(acc, tt.publicId, tt.outputAccount)
			serviceMock := &service.Service{Accounter: acc}
			var brokerMock *broker.Broker

			handler := NewHandler(serviceMock, brokerMock)
			gin.SetMode(gin.ReleaseMode)
			r := gin.New()

			if tt.isTokenExists {
				r.GET("/account/", func(c *gin.Context) {
					c.Set(accountCtx, tt.publicId)
				}, handler.getAccountInfo)
			} else {
				r.GET("/account/", handler.getAccountInfo)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/account/", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedRequestBody, w.Body.String())
		})
	}
}
