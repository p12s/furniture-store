package handler

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// errorResponse - error response
type errorResponse struct {
	Message string `json:"message"`
}

// newErrorResponse - send error
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	if os.Getenv("ENV_CURRENT") == os.Getenv("ENV_PROD") {
		logrus.Printf("%s: [%s] - %s | %s", time.Now().Format(time.RFC3339), c.Request.Method, c.Request.RequestURI, message)
	}
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
