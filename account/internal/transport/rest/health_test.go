package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/p12s/furniture-store/account/internal/broker"
	"github.com/p12s/furniture-store/account/internal/config"
	"github.com/p12s/furniture-store/account/internal/service"
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

func TestHandler_health(t *testing.T) {
	tests := []struct {
		name                string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Can return service name and status",
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"service":"account","status":"OK"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := &service.Service{Accounter: nil}
			var brokerMock *broker.Broker

			handler := NewHandler(serviceMock, brokerMock)
			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.GET("/health", handler.health)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/health", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedRequestBody, w.Body.String())
		})
	}
}
