package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ilya-mezentsev/micro-dep/diagram/internal/clients/shared"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/diagram/stateful"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/diagram/stateless"
	sharedTransport "github.com/ilya-mezentsev/micro-dep/shared/transport/shared"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

type Diagram struct {
	statefulDiagramService  stateful.Service
	statelessDiagramService stateless.Service
}

func NewDiagram(
	statefulDiagramService stateful.Service,
	statelessDiagramService stateless.Service,
) Diagram {

	return Diagram{
		statefulDiagramService:  statefulDiagramService,
		statelessDiagramService: statelessDiagramService,
	}
}

func (d Diagram) Download(context *gin.Context) {
	filepath, err := d.statefulDiagramService.Draw(models.Id(context.Param("id")))
	if err != nil {
		if errors.Is(err, shared.Unauthorized) {
			context.AbortWithStatus(http.StatusUnauthorized)
		} else {
			context.AbortWithStatus(http.StatusInternalServerError)
		}
	} else {
		context.Status(http.StatusOK)
		context.File(filepath)
	}
}

func (d Diagram) Draw(context *gin.Context) {
	rb := sharedTransport.MakeResponseBuilder(context)

	var entities []stateless.Entity
	if err := context.ShouldBindJSON(&entities); err != nil {
		rb.ClientError(err)
		return
	}

	filepath, err := d.statelessDiagramService.Draw(entities)
	if err != nil {
		rb.InternalError(err)
	} else {
		context.Status(http.StatusOK)
		context.File(filepath)
	}
}
