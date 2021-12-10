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
	"github.com/joho/godotenv"
	"github.com/p12s/furniture-store/account/internal/broker"
	mock_broker "github.com/p12s/furniture-store/account/internal/broker/mocks"
	"github.com/p12s/furniture-store/account/internal/config"
	"github.com/p12s/furniture-store/account/internal/domain"
	"github.com/p12s/furniture-store/account/internal/service"
	mock_service "github.com/p12s/furniture-store/account/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

const (
	WAITING_GORUTINE_END_TIME = 1 * time.Second
	DIR_ENV_PATH              = ".env.example"
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

func TestHandler_signUp(t *testing.T) {
	t.Parallel()

	type accountMockBehavior func(s *mock_service.MockAccounter, account domain.Account)
	type brokerMockProducer func(s *mock_broker.MockProducer, event domain.EventType, topic string, input interface{})

	tests := []struct {
		name                string
		inputBody           string
		inputAccount        domain.Account
		accountMockBehavior accountMockBehavior
		eventType           domain.EventType
		topic               string
		eventData           interface{}
		brokerMockProducer  brokerMockProducer
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "Can sign up with correct input",
			inputBody: `{"name": "Ivan", "username": "ivan", "password": "qwerty", "email": "test@test.ru", "address": "Some-city, some-street, some-hause"}`,
			inputAccount: domain.Account{
				Name:     "Ivan",
				Username: "ivan",
				Password: "qwerty",
				Email:    "test@test.ru",
				Address:  "Some-city, some-street, some-hause",
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, account domain.Account) {
				s.EXPECT().CreateAccount(account).Return(nil)
			},
			eventType: domain.EVENT_ACCOUNT_CREATED,
			topic:     "",
			eventData: domain.Account{
				Name:     "Ivan",
				Username: "ivan",
				Password: "",
				Email:    "test@test.ru",
				Address:  "Some-city, some-street, some-hause",
			},
			brokerMockProducer: func(s *mock_broker.MockProducer, event domain.EventType, topic string, input interface{}) {
				s.EXPECT().Produce(event, topic, input).Return(nil)
			},
			expectedStatusCode:  http.StatusCreated,
			expectedRequestBody: `null`,
		},
		{
			name:      "Can't sign up with input without name",
			inputBody: `{"username": "ivan", "password": "qwerty", "email": "test@test.ru", "address": "Some-city, some-street, some-hause"}`,
			inputAccount: domain.Account{
				Username: "ivan",
				Password: "qwerty",
				Email:    "test@test.ru",
				Address:  "Some-city, some-street, some-hause",
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, account domain.Account) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Can't sign up with input without username",
			inputBody: `{"password": "qwerty", "email": "test@test.ru", "address": "Some-city, some-street, some-hause"}`,
			inputAccount: domain.Account{
				Password: "qwerty",
				Email:    "test@test.ru",
				Address:  "Some-city, some-street, some-hause",
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, account domain.Account) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Can't sign up with input without password",
			inputBody: `{"email": "test@test.ru", "address": "Some-city, some-street, some-hause"}`,
			inputAccount: domain.Account{
				Email:   "test@test.ru",
				Address: "Some-city, some-street, some-hause",
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, account domain.Account) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Can't sign up with input without email",
			inputBody: `{"address": "Some-city, some-street, some-hause"}`,
			inputAccount: domain.Account{
				Address: "Some-city, some-street, some-hause",
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, account domain.Account) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:                "Can't sign up with input without address",
			inputBody:           `{}`,
			inputAccount:        domain.Account{},
			accountMockBehavior: func(s *mock_service.MockAccounter, account domain.Account) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Can return error response if service failure",
			inputBody: `{"name": "Ivan", "username": "ivan", "password": "qwerty", "email": "test@test.ru", "address": "Some-city, some-street, some-hause"}`,
			inputAccount: domain.Account{
				Name:     "Ivan",
				Username: "ivan",
				Password: "qwerty",
				Email:    "test@test.ru",
				Address:  "Some-city, some-street, some-hause",
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, account domain.Account) {
				s.EXPECT().CreateAccount(account).Return(errors.New(""))
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
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(tt.inputBody))

			r.ServeHTTP(w, req)
			time.Sleep(WAITING_GORUTINE_END_TIME)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	//t.Parallel()

	type accountMockBehavior func(s *mock_service.MockAccounter, input domain.SignInInput)
	type brokerMockProducer func(s *mock_broker.MockProducer, event domain.EventType, topic string, input interface{})

	tests := []struct {
		name                string
		inputBody           string
		inputAccount        domain.SignInInput
		accountMockBehavior accountMockBehavior
		eventType           domain.EventType
		topic               string
		eventData           interface{}
		brokerMockProducer  brokerMockProducer
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "Can sign in with correct input",
			inputBody: `{"email": "test@test.ru", "password": "qwerty"}`,
			inputAccount: domain.SignInInput{
				Email:    "test@test.ru",
				Password: "qwerty",
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, input domain.SignInInput) {
				s.EXPECT().GenerateTokenByCreds(input.Email, input.Password).Return("token", nil)
			},
			eventType: domain.EVENT_ACCOUNT_TOKEN_UPDATED,
			topic:     "",
			eventData: "token",
			brokerMockProducer: func(s *mock_broker.MockProducer, event domain.EventType, topic string, input interface{}) {
				s.EXPECT().Produce(event, topic, "token").Return(nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"token":"token"}`,
		},
		{
			name:      "Can't sign in with input without email",
			inputBody: `{"password": "qwerty"}`,
			inputAccount: domain.SignInInput{
				Password: "qwerty",
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, input domain.SignInInput) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Can't sign in with input without password",
			inputBody: `{"email": "test@test.ru"}`,
			inputAccount: domain.SignInInput{
				Email: "test@test.ru",
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, input domain.SignInInput) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Can return error response if service failure",
			inputBody: `{"email": "test@test.ru", "password": "qwerty"}`,
			inputAccount: domain.SignInInput{
				Email:    "test@test.ru",
				Password: "qwerty",
			},
			accountMockBehavior: func(s *mock_service.MockAccounter, input domain.SignInInput) {
				s.EXPECT().GenerateTokenByCreds(input.Email, input.Password).Return("", errors.New(""))
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
			r.POST("/sign-in", handler.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(tt.inputBody))

			r.ServeHTTP(w, req)
			time.Sleep(WAITING_GORUTINE_END_TIME)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedRequestBody, w.Body.String())
		})
	}
}
