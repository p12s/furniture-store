package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// errorResponse - error response
type errorResponse struct {
	Message string `json:"message"`
}

// statusResponse - success response
type statusResponse struct {
	Status string `json:"status"`
}

// newErrorResponse - send error
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Printf("%s: [%s] - %s | %s", time.Now().Format(time.RFC3339), c.Request.Method, c.Request.RequestURI, message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
