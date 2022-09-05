package handlers

import (
	log "github.com/angelbirth/logger"

	"auth-poc/svc/auth/application/dto"
	"auth-poc/svc/auth/application/usecase"
	"auth-poc/svc/auth/helper/response"

	"github.com/gin-gonic/gin"
)

type AuthHttpHandler struct {
	AuthUseCase usecase.AuthUseCase
}

func (h *AuthHttpHandler) Login(ctx *gin.Context) {
	resp := response.WrapResponse(ctx)
	form, e := ctx.MultipartForm()
	if e != nil {
		resp.ISE(e.Error(), nil)
		return
	}

	var input = dto.LoginRequest{}

	if v, ok := form.Value["phone_number"]; ok && len(v) > 0 {
		input.PhoneNumber = v[0]
	} else {
		resp.BadRequest("invalid value for field `phone_number`", nil)
		return
	}

	if v, ok := form.Value["password"]; ok && len(v) > 0 {
		input.Password = v[0]
	} else {
		resp.BadRequest("invalid value for field `password`", nil)
		return
	}

	result, err := h.AuthUseCase.Login(&input)
	if err != nil {
		log.Error("[Auth Handler][Login] Err: ", err)
		resp.ISE(err.Error(), nil)
		return
	}

	if result == nil {
		log.Info("[Auth Handler][Login] invalid credentials")
		resp.Unauthorized("invalid credentials")
		return
	}

	resp.OK("Login success!", result)
}

func (h *AuthHttpHandler) RefreshToken(ctx *gin.Context) {
	resp := response.WrapResponse(ctx)

	sentRefreshToken := ctx.Request.Header["Refresh-Token"][0]

	if len(sentRefreshToken) == 0 {
		resp.BadRequest("missing headers", nil)
		return
	}

	result, err := h.AuthUseCase.RefreshToken(sentRefreshToken)
	if err != nil {
		log.Error("[Auth Handler][RefreshToken] Err: ", err)
		resp.ISE(err.Error(), nil)
		return
	}

	resp.OK("Refresh success!", result)
}
