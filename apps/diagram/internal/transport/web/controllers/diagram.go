package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ilya-mezentsev/micro-dep/diagram/internal/clients/shared"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/diagram"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

type Diagram struct {
	diagramService diagram.Service
}

func NewDiagram(diagramService diagram.Service) Diagram {
	return Diagram{diagramService: diagramService}
}

func (d Diagram) Download(context *gin.Context) {
	filepath, err := d.diagramService.Draw(models.Id(context.Param("id")))
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
