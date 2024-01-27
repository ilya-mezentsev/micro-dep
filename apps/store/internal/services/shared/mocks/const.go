package mocks

import (
	"errors"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

var (
	SomeError = errors.New("some-error")

	Relations = []models.Relation{
		{
			Id: "relation-1",
		},
		{
			Id: "relation-2",
		},
		{
			Id: "relation-3",
		},
	}

	Endpoints = []models.Endpoint{
		{
			Id: "endpoint-1",
		},
		{
			Id: "endpoint-2",
		},
		{
			Id: "endpoint-3",
		},
		{
			Id: "endpoint-4",
		},
		{
			Id: "endpoint-5",
		},
		{
			Id: "endpoint-6",
		},
	}

	Entities = []models.Entity{
		{
			Id:        "some-id-1",
			Endpoints: Endpoints[:3],
		},
		{
			Id:        "some-id-2",
			Endpoints: Endpoints[3:],
		},
	}
)
