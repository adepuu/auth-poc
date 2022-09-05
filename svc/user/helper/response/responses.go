package response

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	now = time.Now
)

func WrapResponse(c *gin.Context) *Response {
	return &Response{c}
}

type Response struct {
	*gin.Context
}

func (r *Response) NotFound(message string) {
	r.Context.AbortWithStatusJSON(http.StatusNotFound, BaseResponse{
		Timestamp: now().Unix(),
		Status:    "not_found",
		Message:   message,
	})
}

func (r *Response) NotFoundf(format string, args ...interface{}) {
	r.NotFound(fmt.Sprintf(format, args...))
}

func (r *Response) ISE(message string, data interface{}) {
	r.Context.AbortWithStatusJSON(http.StatusInternalServerError, BaseResponse{
		Timestamp: now().Unix(),
		Status:    "internal_server_error",
		Message:   message,
		Data:      data,
	})
}

func (r *Response) OK(message string, data interface{}) {
	r.Context.JSON(200, BaseResponse{
		Timestamp: now().Unix(),
		Status:    "ok",
		Message:   message,
		Data:      data,
	})
}

func (r *Response) BadRequest(message string, data interface{}) {
	r.Context.AbortWithStatusJSON(http.StatusBadRequest, BaseResponse{
		Timestamp: now().Unix(),
		Status:    "bad_request",
		Message:   message,
		Data:      data,
	})
}

func (r *Response) Created(message string, data interface{}) {
	r.Context.JSON(http.StatusCreated, BaseResponse{
		Timestamp: now().Unix(),
		Status:    "created",
		Message:   message,
		Data:      data,
	})
}
func (r *Response) Conflict(message string) {
	r.Context.JSON(http.StatusConflict, BaseResponse{
		Timestamp: now().Unix(),
		Status:    "conflict",
		Message:   message,
	})
}

func (r *Response) Unauthorized(message string) {
	r.Context.AbortWithStatusJSON(http.StatusUnauthorized, BaseResponse{
		Timestamp: now().Unix(),
		Status:    "unauthorized",
		Message:   message,
	})
}

func (r *Response) Forbidden(message string) {
	r.Context.AbortWithStatusJSON(http.StatusForbidden, BaseResponse{
		Timestamp: now().Unix(),
		Status:    "forbidden",
		Message:   message,
	})
}
