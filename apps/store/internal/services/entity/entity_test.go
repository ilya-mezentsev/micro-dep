package entity

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	entityMocks "github.com/ilya-mezentsev/micro-dep/store/internal/services/entity/mocks"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

var (
	someError = errors.New("some-error")

	allEndpoints = []shared.Endpoint{
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
	entities = []shared.Entity{
		{
			Id:        "some-id-1",
			Endpoints: allEndpoints[:3],
		},
		{
			Id:        "some-id-2",
			Endpoints: allEndpoints[3:],
		},
	}
)

func TestServiceImpl_Create(t *testing.T) {
	tests := []struct {
		name            string
		mockConstructor func() Repo
		expected        error
	}{
		{
			name: "ok",
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().Create(mock.Anything).Return(nil)
				m.EXPECT().Exists(mock.Anything).Return(false, nil)

				return m
			},
			expected: nil,
		},

		{
			name: "failed creation due to repo.Exists error",
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().Exists(mock.Anything).Return(false, someError)

				return m
			},
			expected: someError,
		},

		{
			name: "failed creation due to name existence error",
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().Exists(mock.Anything).Return(true, nil)

				return m
			},
			expected: shared.ExistsError,
		},

		{
			name: "failed creation due to repo.Create error",
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().Create(mock.Anything).Return(someError)
				m.EXPECT().Exists(mock.Anything).Return(false, nil)

				return m
			},
			expected: someError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServiceImpl(tt.mockConstructor())

			require.Equal(t, tt.expected, s.Create(shared.Entity{}))
		})
	}
}

func TestServiceImpl_ReadAll(t *testing.T) {
	tests := []struct {
		name            string
		mockConstructor func() Repo
		expectedDTOs    []shared.Entity
		expectedErr     error
	}{
		{
			name: "ok",
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().ReadAll().Return(entities, nil)

				return m
			},
			expectedDTOs: entities,
			expectedErr:  nil,
		},

		{
			name: "error",
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().ReadAll().Return(nil, someError)

				return m
			},
			expectedDTOs: nil,
			expectedErr:  someError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServiceImpl(tt.mockConstructor())
			dtos, err := s.ReadAll()

			require.Equal(t, tt.expectedDTOs, dtos)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestServiceImpl_ReadOne(t *testing.T) {
	tests := []struct {
		name            string
		entityId        models.Id
		mockConstructor func() Repo
		expectedDTO     shared.Entity
		expectedErr     error
	}{
		{
			name: "ok",
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().ReadOne(entities[0].Id).Return(entities[0], nil)

				return m
			},
			entityId:    entities[0].Id,
			expectedDTO: entities[0],
			expectedErr: nil,
		},

		{
			name: "IdMissingInStorage error",
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().ReadOne(entities[0].Id).Return(shared.Entity{}, errs.IdMissingInStorage)

				return m
			},
			entityId:    entities[0].Id,
			expectedDTO: shared.Entity{},
			expectedErr: shared.NotFoundById,
		},

		{
			name: "general error",
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().ReadOne(entities[0].Id).Return(shared.Entity{}, someError)

				return m
			},
			entityId:    entities[0].Id,
			expectedDTO: shared.Entity{},
			expectedErr: someError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServiceImpl(tt.mockConstructor())
			dtos, err := s.ReadOne(tt.entityId)

			require.Equal(t, tt.expectedDTO, dtos)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestServiceImpl_Update(t *testing.T) {
	tests := []struct {
		name            string
		entityModel     shared.Entity
		mockConstructor func() Repo
		expectedDTO     shared.Entity
		expectedErr     error
	}{
		{
			name:        "ok",
			entityModel: entities[0],
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(entities[0].Id).Return(allEndpoints[:2], nil)
				m.EXPECT().Update(entities[0]).Return(entities[0], nil)

				return m
			},
			expectedDTO: entities[0],
			expectedErr: nil,
		},

		{
			name:        "IdMissingInStorage error",
			entityModel: entities[0],
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(entities[0].Id).Return(nil, errs.IdMissingInStorage)

				return m
			},
			expectedDTO: shared.Entity{},
			expectedErr: shared.NotFoundById,
		},

		{
			name:        "general error",
			entityModel: entities[0],
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(entities[0].Id).Return(nil, someError)

				return m
			},
			expectedDTO: shared.Entity{},
			expectedErr: someError,
		},

		{
			name:        "TryingToRemoveEndpointThatIsInUse error",
			entityModel: entities[0],
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(entities[0].Id).Return(allEndpoints[:4], nil)

				return m
			},
			expectedDTO: shared.Entity{},
			expectedErr: TryingToRemoveEndpointThatIsInUse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServiceImpl(tt.mockConstructor())
			dtos, err := s.Update(tt.entityModel)

			require.Equal(t, tt.expectedDTO, dtos)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestServiceImpl_Delete(t *testing.T) {
	tests := []struct {
		name            string
		entityId        models.Id
		mockConstructor func() Repo
		expectedErr     error
	}{
		{
			name:     "ok",
			entityId: entities[0].Id,
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(entities[0].Id).Return(nil, nil)
				m.EXPECT().Delete(entities[0].Id).Return(nil)

				return m
			},
			expectedErr: nil,
		},

		{
			name:     "IdMissingInStorage error",
			entityId: entities[0].Id,
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(entities[0].Id).Return(nil, errs.IdMissingInStorage)

				return m
			},
			expectedErr: shared.NotFoundById,
		},

		{
			name:     "general error",
			entityId: entities[0].Id,
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(entities[0].Id).Return(nil, someError)

				return m
			},
			expectedErr: someError,
		},

		{
			name:     "TryingToRemoveEntityThatIsUse error",
			entityId: entities[0].Id,
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(entities[0].Id).Return(allEndpoints[:3], nil)

				return m
			},
			expectedErr: TryingToRemoveEntityThatIsUse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServiceImpl(tt.mockConstructor())

			require.Equal(t, tt.expectedErr, s.Delete(tt.entityId))
		})
	}
}
