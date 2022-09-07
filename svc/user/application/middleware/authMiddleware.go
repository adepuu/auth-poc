package middleware

import (
	"auth-poc/svc/user/adapter/grpc/pb"
	"auth-poc/svc/user/constants"
	"auth-poc/svc/user/helper/response"
	"context"
	"strings"

	log "github.com/angelbirth/logger"

	"github.com/gin-gonic/gin"
)

func (i *Interactors) ValidateToken(ctx *gin.Context) {
	rpcCtx, cancel := context.WithTimeout(context.Background(), constants.DEFAULT_TIMEOUT)
	defer cancel()

	resp := response.WrapResponse(ctx)
	authHeader := ctx.GetHeader("Authorization")
	tokenString := authHeader[len(constants.BEARER_SCHEMA):]

	result, err := i.RpcClients.Auth.CheckAuth(rpcCtx, &pb.CheckAuthArgs{
		Token: strings.TrimSpace(tokenString),
	})

	if err != nil {
		log.Error("[Middleware][Auth] Error RPC: ", err)
		resp.ISE("internal server error", nil)
		return
	} else if result.GetIsAuthorized() {
		ctx.Set(constants.USER_ID_CTX, result.GetUserID())
		ctx.Set(constants.USER_TYPE_CTX, int(result.GetUserType()))
		ctx.Next()
		return
	} else {
		log.Infof("[Middleware][Auth] Invalid token '%s'\n", tokenString)
		resp.Unauthorized("invalid or expired credentials")
		return
	}
}

func (i *Interactors) AdminOnly(ctx *gin.Context) {
	resp := response.WrapResponse(ctx)
	userType := ctx.GetInt(constants.USER_TYPE_CTX)

	if userType == int(constants.USER_TYPE_SUPER_ADMIN) {
		ctx.Next()
		return
	} else {
		log.Infof("[Middleware][Auth] Forbidden route")
		resp.Forbidden("you are not supposed to be here")
		return
	}
}
