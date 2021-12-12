package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/p12s/furniture-store/account/internal/domain"
	"github.com/sirupsen/logrus"
)

// @Summary Update account
// @Tags Account
// @Description Update account info
// @ID updateAccount
// @Accept  json
// @Param input body domain.UpdateAccountInput true "credentials"
// @Success 200
// @Router /account/info [put]
func (h *Handler) updateAccount(c *gin.Context) { // nolint
	var input domain.UpdateAccountInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.services.UpdateAccountInfo(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "service failure")
		return
	}

	go func() {
		err := h.broker.Produce(domain.EVENT_ACCOUNT_INFO_UPDATED, h.broker.TopicAccountCUD, input)
		if err != nil {
			logrus.Errorf("sent update account event fail: %s/n", err.Error())
		}
	}()

	c.Status(http.StatusOK)
}

// @Summary Delete account
// @Tags Account
// @Description Delete account
// @ID deleteAccount
// @Accept  json
// @Param input body domain.DeleteAccountInput true "credentials"
// @Success 200
// @Router /account/ [delete]
func (h *Handler) deleteAccount(c *gin.Context) {
	var input domain.DeleteAccountInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.services.DeleteAccount(input.PublicId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "service failure")
		return
	}

	go func() {
		err := h.broker.Produce(domain.EVENT_ACCOUNT_DELETED, h.broker.TopicAccountCUD, input)
		if err != nil {
			logrus.Errorf("sent delete account event fail: %s/n", err.Error())
		}
	}()

	c.Status(http.StatusOK)
}

// @Summary Get account info
// @Tags Account
// @Description Get account
// @ID getAccountInfo
// @Accept  json
// @Success 200
// @Router /account/ [get]
func (h *Handler) getAccountInfo(c *gin.Context) {
	accountPublicId, err := getAccountPublicId(c)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "account public id not found")
		return
	}

	account, err := h.services.GetAccount(accountPublicId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "service failure")
		return
	}

	c.JSON(http.StatusOK, account)
}
