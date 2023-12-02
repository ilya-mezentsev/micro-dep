package endpoint

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	endpointMocks "github.com/ilya-mezentsev/micro-dep/store/internal/services/endpoint/mocks"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
	sharedMocks "github.com/ilya-mezentsev/micro-dep/store/internal/services/shared/mocks"
)

func TestServiceImpl_Create(t *testing.T) {
	tests := []struct {
		name            string
		model           shared.Endpoint
		mockConstructor func() Repo
		expected        error
	}{
		{
			name:  "ok",
			model: sharedMocks.Endpoints[0],
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().Exists(sharedMocks.Endpoints[0]).Return(true, false, nil)
				m.EXPECT().Create(sharedMocks.Endpoints[0]).Return(shared.Endpoint{}, nil)

				return m
			},
			expected: nil,
		},

		{
			name:  "IdMissingInStorage error",
			model: sharedMocks.Endpoints[0],
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().Exists(sharedMocks.Endpoints[0]).Return(false, false, errs.IdMissingInStorage)

				return m
			},
			expected: shared.NotFoundById,
		},

		{
			name:  "some error",
			model: sharedMocks.Endpoints[0],
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().Exists(sharedMocks.Endpoints[0]).Return(false, false, sharedMocks.SomeError)

				return m
			},
			expected: sharedMocks.SomeError,
		},

		{
			name:  "entity does not exists error",
			model: sharedMocks.Endpoints[0],
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().Exists(sharedMocks.Endpoints[0]).Return(false, false, nil)

				return m
			},
			expected: TryingToAddEndpointToMissingEntity,
		},

		{
			name:  "endpoint exists error",
			model: sharedMocks.Endpoints[0],
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().Exists(sharedMocks.Endpoints[0]).Return(true, true, nil)

				return m
			},
			expected: TryingToCreateEndpointThatExists,
		},

		{
			name:  "creation error",
			model: sharedMocks.Endpoints[0],
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().Exists(sharedMocks.Endpoints[0]).Return(true, false, nil)
				m.EXPECT().Create(sharedMocks.Endpoints[0]).Return(shared.Endpoint{}, sharedMocks.SomeError)

				return m
			},
			expected: sharedMocks.SomeError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServiceImpl(tt.mockConstructor())
			_, err := s.Create(tt.model)

			require.Equal(t, tt.expected, err)
		})
	}
}

func TestServiceImpl_Update(t *testing.T) {
	tests := []struct {
		name            string
		model           shared.Endpoint
		mockConstructor func() Repo
		expectedModel   shared.Endpoint
		expectedError   error
	}{
		{
			name:  "ok",
			model: sharedMocks.Endpoints[0],
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().Exists(sharedMocks.Endpoints[0]).Return(true, true, nil)
				m.EXPECT().Update(sharedMocks.Endpoints[0]).Return(sharedMocks.Endpoints[0], nil)

				return m
			},
			expectedModel: sharedMocks.Endpoints[0],
			expectedError: nil,
		},

		{
			name:  "IdMissingInStorage error",
			model: sharedMocks.Endpoints[0],
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().Exists(sharedMocks.Endpoints[0]).Return(false, false, errs.IdMissingInStorage)

				return m
			},
			expectedModel: shared.Endpoint{},
			expectedError: shared.NotFoundById,
		},

		{
			name:  "some error",
			model: sharedMocks.Endpoints[0],
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().Exists(sharedMocks.Endpoints[0]).Return(false, false, sharedMocks.SomeError)

				return m
			},
			expectedModel: shared.Endpoint{},
			expectedError: sharedMocks.SomeError,
		},

		{
			name:  "endpoint does not exists error",
			model: sharedMocks.Endpoints[0],
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().Exists(sharedMocks.Endpoints[0]).Return(true, false, nil)

				return m
			},
			expectedModel: shared.Endpoint{},
			expectedError: TryingToUpdateMissingEndpoint,
		},

		{
			name:  "update error",
			model: sharedMocks.Endpoints[0],
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().Exists(sharedMocks.Endpoints[0]).Return(true, true, nil)
				m.EXPECT().Update(sharedMocks.Endpoints[0]).Return(shared.Endpoint{}, sharedMocks.SomeError)

				return m
			},
			expectedModel: shared.Endpoint{},
			expectedError: sharedMocks.SomeError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServiceImpl(tt.mockConstructor())
			model, err := s.Update(tt.model)

			require.Equal(t, tt.expectedModel, model)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestServiceImpl_Delete(t *testing.T) {
	tests := []struct {
		name            string
		modelId         models.Id
		mockConstructor func() Repo
		expected        error
	}{
		{
			name:    "ok",
			modelId: sharedMocks.Endpoints[0].Id,
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().HasRelation(sharedMocks.Endpoints[0].Id).Return(false, nil)
				m.EXPECT().Delete(sharedMocks.Endpoints[0].Id).Return(nil)

				return m
			},
			expected: nil,
		},

		{
			name:    "IdMissingInStorage error",
			modelId: sharedMocks.Endpoints[0].Id,
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().HasRelation(sharedMocks.Endpoints[0].Id).Return(false, errs.IdMissingInStorage)

				return m
			},
			expected: shared.NotFoundById,
		},

		{
			name:    "some error",
			modelId: sharedMocks.Endpoints[0].Id,
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().HasRelation(sharedMocks.Endpoints[0].Id).Return(false, sharedMocks.SomeError)

				return m
			},
			expected: sharedMocks.SomeError,
		},

		{
			name:    "relation exists error",
			modelId: sharedMocks.Endpoints[0].Id,
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().HasRelation(sharedMocks.Endpoints[0].Id).Return(true, nil)

				return m
			},
			expected: TryingToRemoveEndpointThatHasRelation,
		},

		{
			name:    "deletion error",
			modelId: sharedMocks.Endpoints[0].Id,
			mockConstructor: func() Repo {
				m := endpointMocks.NewMockRepo(t)
				m.EXPECT().HasRelation(sharedMocks.Endpoints[0].Id).Return(false, nil)
				m.EXPECT().Delete(sharedMocks.Endpoints[0].Id).Return(sharedMocks.SomeError)

				return m
			},
			expected: sharedMocks.SomeError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServiceImpl(tt.mockConstructor())

			require.Equal(t, tt.expected, s.Delete(tt.modelId))
		})
	}
}
