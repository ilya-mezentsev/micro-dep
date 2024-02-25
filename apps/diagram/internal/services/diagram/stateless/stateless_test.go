package stateless

import (
	"errors"
	"io"
	"log/slog"
	"sort"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/diagram/shared"
	sharedMocks "github.com/ilya-mezentsev/micro-dep/diagram/internal/services/diagram/shared/mocks"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/shared/types"
	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

var (
	mockId    = "mock-id"
	someError = errors.New("some-error")
)

func TestService_Draw(t *testing.T) {
	tests := []struct {
		name            string
		entities        []Entity
		mockConstructor func() shared.DrawService
		expectedResult  string
		expectedError   error
	}{
		{
			name: "ok",
			mockConstructor: func() shared.DrawService {
				ds := sharedMocks.NewMockDrawService(t)
				ds.EXPECT().DrawDiagram(mock.Anything).Return("path", nil)

				return ds
			},
			expectedResult: "path",
		},

		{
			name: "error",
			mockConstructor: func() shared.DrawService {
				ds := sharedMocks.NewMockDrawService(t)
				ds.EXPECT().DrawDiagram(mock.Anything).Return("", someError)

				return ds
			},
			expectedError: errs.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newServiceWithMocks(tt.mockConstructor(), mockId)

			result, err := s.Draw(tt.entities)

			require.Equal(t, tt.expectedResult, result)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestService_BuildRelationsDiagramData(t *testing.T) {
	s := newServiceWithMocks(nil, mockId)

	entities := []Entity{
		{
			Name: "s-1",
			Dependencies: []Entity{
				{
					Name: "s-2",
					Endpoints: []Endpoint{
						{
							Kind:    "kind-2",
							Address: "address-2",
						},
					},
				},
				{
					Name: "s-3",
					Endpoints: []Endpoint{
						{
							Kind:    "kind-3",
							Address: "address-3",
						},
					},
				},
			},
		},

		{
			Name: "s-2",
			Dependencies: []Entity{
				{
					Name: "s-3",
					Endpoints: []Endpoint{
						{
							Kind:    "kind-3",
							Address: "address-31",
						},
					},
				},
			},
		},
	}

	rdd := types.RelationsDiagramData{
		Entities: []models.Entity{
			{
				Id:        models.Id(mockId),
				Name:      "s-1",
				Endpoints: make([]models.Endpoint, 0),
			},
			{
				Id:   models.Id(mockId),
				Name: "s-2",
				Endpoints: []models.Endpoint{
					{
						Id:       models.Id(mockId),
						EntityId: models.Id(mockId),
						Kind:     "kind-2",
						Address:  "address-2",
					},
				},
			},
			{
				Id:   models.Id(mockId),
				Name: "s-3",
				Endpoints: []models.Endpoint{
					{
						Id:       models.Id(mockId),
						EntityId: models.Id(mockId),
						Kind:     "kind-3",
						Address:  "address-3",
					},
					{
						Id:       models.Id(mockId),
						EntityId: models.Id(mockId),
						Kind:     "kind-3",
						Address:  "address-31",
					},
				},
			},
		},

		// fixme: cannot check ids correctness
		Relations: []models.Relation{
			{
				Id:           models.Id(mockId),
				FromEntityId: models.Id(mockId),
				ToEndpointId: models.Id(mockId),
			},
			{
				Id:           models.Id(mockId),
				FromEntityId: models.Id(mockId),
				ToEndpointId: models.Id(mockId),
			},
			{
				Id:           models.Id(mockId),
				FromEntityId: models.Id(mockId),
				ToEndpointId: models.Id(mockId),
			},
		},
	}

	result := s.buildRelationsDiagramData(entities)

	// sort entities to avoid order issues (result is built from map)
	sort.Slice(result.Entities, func(i, j int) bool {
		return result.Entities[i].Name < result.Entities[j].Name
	})

	require.Equal(t, rdd, result)
}

func TestService_BuildEntityModel(t *testing.T) {
	s := newServiceWithMocks(nil, mockId)

	tests := []struct {
		name     string
		entity   Entity
		expected models.Entity
	}{
		{
			name: "build entity without endpoints",
			entity: Entity{
				Name: "e-1",
			},
			expected: models.Entity{
				Id:        models.Id(mockId),
				Name:      "e-1",
				Endpoints: make([]models.Endpoint, 0),
			},
		},

		{
			name: "build entity without endpoints",
			entity: Entity{
				Name: "e-1",
				Endpoints: []Endpoint{
					{
						Kind:    "kind-1",
						Address: "address-1",
					},
					{
						Kind:    "kind-2",
						Address: "address-1",
					},
				},
			},
			expected: models.Entity{
				Id:   models.Id(mockId),
				Name: "e-1",
				Endpoints: []models.Endpoint{
					{
						Id:       models.Id(mockId),
						EntityId: models.Id(mockId),
						Kind:     "kind-1",
						Address:  "address-1",
					},
					{
						Id:       models.Id(mockId),
						EntityId: models.Id(mockId),
						Kind:     "kind-2",
						Address:  "address-1",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := s.buildEntityModel(tt.entity)

			require.Equal(t, tt.expected, entity)
		})
	}

	// just for coverage
	s = New(nil, slogMock())
	e := s.buildEntityModel(tests[0].entity)
	require.Equal(t, 36, len(e.Id))
}

func TestService_MergeEndpoints(t *testing.T) {
	s := newServiceWithMocks(nil, mockId)

	tests := []struct {
		name      string
		entity    models.Entity
		endpoints []Endpoint
		expected  models.Entity
	}{
		{
			name: "merge one endpoint",
			entity: models.Entity{
				Id: models.Id(mockId),
			},
			endpoints: []Endpoint{
				{
					Kind:    "kind-1",
					Address: "address-1",
				},
			},
			expected: models.Entity{
				Id: models.Id(mockId),
				Endpoints: []models.Endpoint{
					{
						Id:       models.Id(mockId),
						EntityId: models.Id(mockId),
						Kind:     "kind-1",
						Address:  "address-1",
					},
				},
			},
		},

		{
			name: "merge one new endpoint",
			entity: models.Entity{
				Id: models.Id(mockId),
				Endpoints: []models.Endpoint{
					{
						Id:       models.Id(mockId),
						EntityId: models.Id(mockId),
						Kind:     "kind-1",
						Address:  "address-1",
					},
				},
			},
			endpoints: []Endpoint{
				{
					Kind:    "kind-2",
					Address: "address-1",
				},
			},
			expected: models.Entity{
				Id: models.Id(mockId),
				Endpoints: []models.Endpoint{
					{
						Id:       models.Id(mockId),
						EntityId: models.Id(mockId),
						Kind:     "kind-1",
						Address:  "address-1",
					},
					{
						Id:       models.Id(mockId),
						EntityId: models.Id(mockId),
						Kind:     "kind-2",
						Address:  "address-1",
					},
				},
			},
		},

		{
			name: "skip exists endpoint",
			entity: models.Entity{
				Id: models.Id(mockId),
				Endpoints: []models.Endpoint{
					{
						Id:       models.Id(mockId),
						EntityId: models.Id(mockId),
						Kind:     "kind-1",
						Address:  "address-1",
					},
				},
			},
			endpoints: []Endpoint{
				{
					Kind:    "kind-1",
					Address: "address-1",
				},
				{
					Kind:    "kind-2",
					Address: "address-1",
				},
			},
			expected: models.Entity{
				Id: models.Id(mockId),
				Endpoints: []models.Endpoint{
					{
						Id:       models.Id(mockId),
						EntityId: models.Id(mockId),
						Kind:     "kind-1",
						Address:  "address-1",
					},
					{
						Id:       models.Id(mockId),
						EntityId: models.Id(mockId),
						Kind:     "kind-2",
						Address:  "address-1",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := s.mergeEndpoints(tt.entity, tt.endpoints)

			require.Equal(t, tt.expected, entity)
		})
	}
}

func TestService_FilterEndpoints(t *testing.T) {
	s := newServiceWithMocks(nil, mockId)

	tests := []struct {
		name      string
		entity    models.Entity
		endpoints []Endpoint
		expected  []models.Endpoint
	}{
		{
			name: "filter single exists endpoint",
			entity: models.Entity{
				Endpoints: []models.Endpoint{
					{
						Kind:    "kind-1",
						Address: "address-1",
					},
					{
						Kind:    "kind-2",
						Address: "address-1",
					},
				},
			},
			endpoints: []Endpoint{
				{
					Kind:    "kind-2",
					Address: "address-1",
				},
			},
			expected: []models.Endpoint{
				{
					Kind:    "kind-2",
					Address: "address-1",
				},
			},
		},

		{
			name: "filter exists and not exists endpoint",
			entity: models.Entity{
				Endpoints: []models.Endpoint{
					{
						Kind:    "kind-1",
						Address: "address-1",
					},
					{
						Kind:    "kind-2",
						Address: "address-1",
					},
				},
			},
			endpoints: []Endpoint{
				{
					Kind:    "kind-2",
					Address: "address-1",
				},
				{
					Kind:    "kind-2",
					Address: "address-2",
				},
			},
			expected: []models.Endpoint{
				{
					Kind:    "kind-2",
					Address: "address-1",
				},
			},
		},

		{
			name: "filter single not exists endpoint",
			entity: models.Entity{
				Endpoints: []models.Endpoint{
					{
						Kind:    "kind-1",
						Address: "address-1",
					},
					{
						Kind:    "kind-2",
						Address: "address-1",
					},
				},
			},
			endpoints: []Endpoint{
				{
					Kind:    "kind-2",
					Address: "address-2",
				},
			},
			expected: []models.Endpoint{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := s.filterEndpoints(tt.entity, tt.endpoints)
			require.Equal(t, tt.expected, res)
		})
	}
}

func idFactoryMock(id string) func() models.Id {
	return func() models.Id {
		return models.Id(id)
	}
}

func slogMock() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newServiceWithMocks(drawService shared.DrawService, id string) Service {
	return Service{
		drawService: drawService,
		logger:      slogMock(),
		idFactory:   idFactoryMock(id),
	}
}
