package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/p12s/furniture-store/account/internal/broker"
	mock_broker "github.com/p12s/furniture-store/account/internal/broker/mocks"
	"github.com/p12s/furniture-store/account/internal/domain"
	"github.com/p12s/furniture-store/account/internal/service"
	mock_service "github.com/p12s/furniture-store/account/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_userIdentity(t *testing.T) {
	//t.Parallel()

	type accountMockBehavior func(s *mock_service.MockAccounter, token string)
	type brokerMockProducer func(s *mock_broker.MockProducer, event domain.EventType, topic string, input interface{})

	tests := []struct {
		name                string
		headerName          string
		headerValue         string
		token               string
		accountMockBehavior accountMockBehavior
		expectedStatusCode  int
		expectedRequestBody string
		brokerMockProducer  brokerMockProducer
	}{ /*
			{
				name:        "Can identify right token",
				headerName:  "Authorization",
				headerValue: "Bearer token",
				token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzcyNzgxMjQsImlhdCI6MTYzNzIzNDkyNCwiYWNjb3VudF9pZCI6NH0.QiTQv3yHwYQwdSKQ3FwFFoMnlq07lSQwYm2w4tozfA0",
				accountMockBehavior: func(s *mock_service.MockAccounter, token string) {
					s.EXPECT().ParseToken(token).Return("1", nil)
				},
				expectedStatusCode:  http.StatusOK,
				expectedRequestBody: "1",
			},*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			acc := mock_service.NewMockAccounter(ctrl)
			tt.accountMockBehavior(acc, tt.token)
			serviceMock := &service.Service{Accounter: acc}
			var brokerMock *broker.Broker

			handler := NewHandler(serviceMock, brokerMock)
			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.POST("/account", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get(accountCtx)
				c.String(http.StatusOK, fmt.Sprintf("%d", id.(int)))
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/account", nil)
			req.Header.Set(tt.headerName, tt.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)           // 404
			assert.Equal(t, tt.expectedRequestBody, w.Body.String()) // 404 page not found
		})
	}
}
