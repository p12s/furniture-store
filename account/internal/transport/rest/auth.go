package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/p12s/furniture-store/account/internal/domain"
	"github.com/sirupsen/logrus"
)

// @Summary Sign up
// @Tags Auth
// @Description Create account
// @ID signUp
// @Accept  json
// @Param input body domain.Account true "credentials"
// @Success 201
// @Router /sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input domain.Account
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.services.CreateAccount(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "service failure")
		return
	}

	input.Password = ""
	go func() {
		err := h.broker.Produce(domain.EVENT_ACCOUNT_CREATED, h.broker.TopicAccountCUD, input)
		if err != nil {
			logrus.Errorf("sent sign-up event fail: %s/n", err.Error())
		}
	}()

	c.Status(http.StatusCreated)
}

// @Summary Sign in
// @Tags Auth
// @Description Sending data to get authentication with jwt-token
// @ID signIn
// @Accept  json
// @Produce  json
// @Param input body domain.SignInInput true "credentials"
// @Success 200
// @Router /sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input domain.SignInInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	accountToken, err := h.services.Accounter.GenerateTokenByCreds(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "service failure")
		return
	}

	go func() {
		err := h.broker.Producer.Produce(domain.EVENT_ACCOUNT_TOKEN_UPDATED, h.broker.TopicAccountCUD, accountToken)
		if err != nil {
			logrus.Errorf("sent sign-in event fail: %s/n", err.Error())
		}
	}()

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": accountToken,
	})
}
