package shared

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	ResponseBuilder struct {
		context  *gin.Context
		response any
	}

	okResponse struct {
		Data any `json:"data"`
	}

	errorResponse struct {
		Error string `json:"error"`
	}
)

func MakeResponseBuilder(context *gin.Context) ResponseBuilder {
	return ResponseBuilder{
		context: context,
	}
}

func (r ResponseBuilder) InternalError(err error) {
	r.context.JSON(http.StatusInternalServerError, errorResponse{Error: cutError(err)})
}

func (r ResponseBuilder) NotImplemented() {
	r.context.Status(http.StatusNotImplemented)
}

func (r ResponseBuilder) NotFoundError(err error) {
	r.context.JSON(http.StatusNotFound, errorResponse{Error: cutError(err)})
}

func (r ResponseBuilder) EmptyNotFound() {
	r.context.AbortWithStatus(http.StatusNotFound)
}

func (r ResponseBuilder) ClientError(err error) {
	r.context.JSON(http.StatusBadRequest, errorResponse{Error: cutError(err)})
}

func (r ResponseBuilder) UnauthorizedError(err error) {
	r.context.JSON(http.StatusUnauthorized, errorResponse{Error: cutError(err)})
}

func (r ResponseBuilder) ConflictError(err error) {
	r.context.JSON(http.StatusConflict, errorResponse{Error: cutError(err)})
}

func (r ResponseBuilder) Ok(response any) {
	r.context.JSON(http.StatusOK, okResponse{Data: response})
}

func (r ResponseBuilder) Created() {
	r.context.Status(http.StatusCreated)
}

func (r ResponseBuilder) EmptyOk() {
	r.context.Status(http.StatusNoContent)
}

func cutError(err error) string {
	msg := err.Error()

	return strings.Split(msg, "\n")[0]
}
