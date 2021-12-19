package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/p12s/furniture-store/product/internal/broker"
	"github.com/p12s/furniture-store/product/internal/service"
)

type Handler struct {
	services *service.Service
	broker   *broker.Broker
}

func NewHandler(services *service.Service, broker *broker.Broker) *Handler {
	return &Handler{services: services, broker: broker}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(CORSMiddleware())

	router.GET("/health", h.health)
	router.POST("/sign-up", h.signUp)
	router.POST("/sign-in", h.signIn)

	account := router.Group("/account", h.userIdentity)
	{
		// account.GET("/", h.getAccountInfo)
		// account.PUT("/info", h.updateAccount)
		// account.PUT("/role", h.updateAccountRole)
		account.DELETE("/", h.deleteAccount)
	}

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH,OPTIONS,GET,PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
