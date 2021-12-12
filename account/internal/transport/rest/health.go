package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary Health
// @Tags Service
// @Description Health check
// @ID health
// @Success 200
// @Router /health [get]
func (h *Handler) health(c *gin.Context) {
	if os.Getenv("ENV_CURRENT") == os.Getenv("ENV_PROD") {
		logrus.Printf("%s: [%s] - %s ", time.Now().Format(time.RFC3339), c.Request.Method, c.Request.RequestURI)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"service": "account",
		"status":  "OK",
	})
}
