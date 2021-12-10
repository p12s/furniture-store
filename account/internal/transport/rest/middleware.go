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

	if headerParts[0] != "Bearer" {
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

func getAccountPublicId(c *gin.Context) (string, error) { // nolint
	value, ok := c.Get(accountCtx)

	if !ok {
		return "", errors.New("account publicId not found")
	}

	publicId, ok := value.(string)
	if !ok {
		return "", errors.New("account publicId is of invalid type")
	}

	return publicId, nil
}
