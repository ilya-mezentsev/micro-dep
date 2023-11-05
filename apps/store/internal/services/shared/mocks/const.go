package mocks

import (
	"errors"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

var (
	SomeError = errors.New("some-error")

	AllEndpoints = []shared.Endpoint{
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

	Entities = []shared.Entity{
		{
			Id:        "some-id-1",
			Endpoints: AllEndpoints[:3],
		},
		{
			Id:        "some-id-2",
			Endpoints: AllEndpoints[3:],
		},
	}
)
