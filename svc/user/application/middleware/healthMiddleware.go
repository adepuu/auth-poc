package middleware

import (
	"auth-poc/svc/user/constants"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (i *Interactors) HealthCheck(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), constants.DEFAULT_TIMEOUT)
	defer cancel()

	err := i.DataStore.Client.Ping(c, nil)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	ctx.AbortWithStatus(http.StatusOK)
}
