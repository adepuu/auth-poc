package handlers

import (
	"auth-poc/svc/user/application/dto"
	"auth-poc/svc/user/application/usecase"
	"auth-poc/svc/user/constants"
	"auth-poc/svc/user/helper/response"
	"strconv"

	log "github.com/angelbirth/logger"
	"github.com/gin-gonic/gin"
)

type UserHttpHandler struct {
	UserUseCase usecase.UserUseCase
}

func (h *UserHttpHandler) Register(ctx *gin.Context) {
	resp := response.WrapResponse(ctx)
	asAdminQuestion := ctx.DefaultQuery("as_admin", "")

	var input dto.RegisterRequest

	err := ctx.BindJSON(&input)
	if err != nil {
		resp.BadRequest("invalid json data", err)
		return
	}

	// Testing only -> handle admin register
	// TODO: Remove after creating init DB data hooks
	if input.UserType == uint32(constants.USER_TYPE_SUPER_ADMIN) && asAdminQuestion != constants.AS_ADMIN_VALUE {
		log.Infof("[User Handler][Register] Cant register as Admin")
		resp.Forbidden("Cant register as Admin: it's dangerous")
		return
	}

	if input.UserType == 0 {
		input.UserType = uint32(constants.USER_TYPE_REGULAR)
	}

	if len(input.PhoneNumber) <= 6 {
		resp.BadRequest("phone number too short", nil)
		return
	}

	if len(input.FullName) == 0 || len(input.Email) == 0 {
		resp.BadRequest("full name and email field is mandatory", nil)
		return
	}

	result, err := h.UserUseCase.Register(&input)
	if err != nil {
		log.Error("[User Handler][Register] Err: ", err)
		resp.ISE(err.Error(), nil)
		return
	}
	if result == nil {
		resp.Conflict("phone number already used by other user")
		return
	}

	resp.Created("Register success!", result)
}

func (h *UserHttpHandler) Profile(ctx *gin.Context) {
	resp := response.WrapResponse(ctx)
	userID := ctx.GetString(constants.USER_ID_CTX)
	if len(userID) == 0 {
		resp.BadRequest("malformed data", nil)
		return
	}

	userData, err := h.UserUseCase.GetUserByKey(userID, "")
	if err != nil {
		log.Error("[User Handler][Profile] Err: ", err)
		resp.ISE(err.Error(), nil)
		return
	}
	resp.OK("Get profile success!", userData)
}

func (h *UserHttpHandler) Detail(ctx *gin.Context) {
	resp := response.WrapResponse(ctx)
	userID := ctx.Param("id")
	if len(userID) == 0 {
		resp.BadRequest("malformed data", nil)
		return
	}

	userData, err := h.UserUseCase.GetUserByKey(userID, "")
	if err != nil {
		log.Error("[User Handler][Detail] Err: ", err)
		resp.ISE(err.Error(), nil)
		return
	}
	resp.OK("Get detail success!", userData)
}

func (h *UserHttpHandler) AllUser(ctx *gin.Context) {
	resp := response.WrapResponse(ctx)

	page, err := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		resp.BadRequest("invalid page data", err)
		return
	}
	limit, err := strconv.ParseInt(ctx.DefaultQuery("limit", "1"), 10, 64)
	if err != nil {
		resp.BadRequest("invalid limit data", err)
		return
	}

	userData, err := h.UserUseCase.GetAllUser(limit, page)
	if err != nil {
		log.Error("[User Handler][Detail] Err: ", err)
		resp.ISE(err.Error(), nil)
		return
	}
	resp.OK("Get detail success!", userData)
}

func (h *UserHttpHandler) Delete(ctx *gin.Context) {
	resp := response.WrapResponse(ctx)
	userID := ctx.Param("id")
	if len(userID) == 0 {
		resp.BadRequest("malformed data", nil)
		return
	}

	requesterID := ctx.GetString(constants.USER_ID_CTX)
	if userID == requesterID {
		resp.BadRequest("can't delete self", nil)
		return
	}

	err := h.UserUseCase.DeleteUser(userID)
	if err != nil {
		log.Error("[User Handler][Delete] Err: ", err)
		resp.ISE(err.Error(), nil)
		return
	}
	resp.OK("User deleted!", nil)
}

func (h *UserHttpHandler) Update(ctx *gin.Context) {
	resp := response.WrapResponse(ctx)

	var input dto.UpdateRequest
	err := ctx.BindJSON(&input)
	if err != nil {
		resp.BadRequest("invalid json data", err)
		return
	}

	input.UserID = ctx.Param("id")
	if len(input.UserID) == 0 {
		resp.BadRequest("malformed data", nil)
		return
	}

	if len(input.FullName) == 0 || len(input.Email) == 0 {
		resp.BadRequest("full name and email field is mandatory", nil)
		return
	}

	result, err := h.UserUseCase.Update(&input)
	if err != nil {
		log.Error("[User Handler][Update] Err: ", err)
		resp.ISE(err.Error(), nil)
		return
	}

	resp.OK("User updated!", result)
}
