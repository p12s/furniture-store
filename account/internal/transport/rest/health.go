package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) health(c *gin.Context) {
	logrus.Printf("%s: [%s] - %s ", time.Now().Format(time.RFC3339), c.Request.Method, c.Request.RequestURI)

	c.JSON(http.StatusOK, map[string]interface{}{
		"service": "account",
		"status":  "OK",
	})
}
