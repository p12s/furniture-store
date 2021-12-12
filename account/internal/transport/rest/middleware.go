package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHandler = "Authorization"
	accountCtx           = "accountPublicId"
)

// userIdentity - checking token
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHandler)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if headerParts[1] == "" {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	accountId, err := h.services.Accounter.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	c.Set(accountCtx, accountId)
}

// getAccountPublicId - getting current account public_id
func getAccountPublicId(c *gin.Context) (string, error) {
	id, ok := c.Get(accountCtx)
	if !ok {
		return "", errors.New("account public_id not found")
	}

	idString, ok := id.(string)
	if !ok {
		return "", errors.New("account id is of invalid type")
	}

	return idString, nil
}
