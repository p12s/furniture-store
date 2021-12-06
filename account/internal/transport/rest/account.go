package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/p12s/furniture-store/account/internal/domain"
)

func (h *Handler) updateAccount(c *gin.Context) { // nolint
	var input domain.UpdateAccountInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.services.UpdateAccountInfo(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	go h.broker.Event(domain.EVENT_ACCOUNT_INFO_UPDATED, "h.broker.TopicAccountBE", input)

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "OK",
	})
}

func (h *Handler) deleteAccount(c *gin.Context) {
	// TODO надо проверять что это 1) либо свой акк - свой могу удалить всегда 2) либо есть роль админа
	// иначе - нет прав
	var input domain.DeleteAccountInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.services.DeleteAccount(input.PublicId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	go h.broker.Event(domain.EVENT_ACCOUNT_DELETED, "h.broker.TopicAccountCUD", input)

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "OK",
	})
}
