package shared

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	ResponseBuilder interface {
		EmptyOk()
		Created()
		Ok(response any)
		NotImplemented()
		InternalError(err error)
		ClientError(err error)
		UnauthorizedError(err error)
		ConflictError(err error)
		NotFoundError(err error)
	}

	responseBuilder struct {
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
	return responseBuilder{
		context: context,
	}
}

func (r responseBuilder) InternalError(err error) {
	r.context.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
}

func (r responseBuilder) NotImplemented() {
	r.context.Status(http.StatusNotImplemented)
}

func (r responseBuilder) NotFoundError(err error) {
	r.context.JSON(http.StatusNotFound, errorResponse{Error: err.Error()})
}

func (r responseBuilder) ClientError(err error) {
	r.context.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
}

func (r responseBuilder) UnauthorizedError(err error) {
	r.context.JSON(http.StatusUnauthorized, errorResponse{Error: err.Error()})
}

func (r responseBuilder) ConflictError(err error) {
	r.context.JSON(http.StatusConflict, errorResponse{Error: err.Error()})
}

func (r responseBuilder) Ok(response any) {
	r.context.JSON(http.StatusOK, okResponse{Data: response})
}

func (r responseBuilder) Created() {
	r.context.Status(http.StatusCreated)
}

func (r responseBuilder) EmptyOk() {
	r.context.Status(http.StatusNoContent)
}
