package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/p12s/furniture-store/account/internal/domain"
	"github.com/sirupsen/logrus"
)

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

	c.JSON(http.StatusCreated, nil)
}

// TODO здесь и в других сервисах в токене юзера храню account_id (int) -
// первичный ключ из БД сервиса auth
// по-хорошему надо хранить public_id (uuid.UUID) - он одинаковый во всех сервисах
// в то время как account_id (int) может отличаться
func (h *Handler) signIn(c *gin.Context) {
	var input domain.SignInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accountToken, err := h.services.Accounter.GenerateTokenByCreds(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	go h.broker.Producer.Produce(domain.EVENT_ACCOUNT_TOKEN_UPDATED, "h.broker.TopicAccountCUD", accountToken)

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": accountToken,
	})
}

func (h *Handler) token(c *gin.Context) { // nolint
	accountId, err := getAccountId(c)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	_ = accountId
	/*
		account, err := h.services.GetAccountById(accountId) // TODO доабвить GetAccountById
		// TODO если "no rows in result set" - возвращать осмысленный текст
		if err != nil {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
	*/

	c.JSON(http.StatusOK, nil)
}
