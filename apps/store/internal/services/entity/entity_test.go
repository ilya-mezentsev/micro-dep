package entity

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	entityMocks "github.com/ilya-mezentsev/micro-dep/store/internal/services/entity/mocks"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
	sharedMocks "github.com/ilya-mezentsev/micro-dep/store/internal/services/shared/mocks"
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
				m.EXPECT().Exists(mock.Anything).Return(false, sharedMocks.SomeError)

				return m
			},
			expected: sharedMocks.SomeError,
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
				m.EXPECT().Create(mock.Anything).Return(sharedMocks.SomeError)
				m.EXPECT().Exists(mock.Anything).Return(false, nil)

				return m
			},
			expected: sharedMocks.SomeError,
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
				m.EXPECT().ReadAll().Return(sharedMocks.Entities, nil)

				return m
			},
			expectedDTOs: sharedMocks.Entities,
			expectedErr:  nil,
		},

		{
			name: "error",
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().ReadAll().Return(nil, sharedMocks.SomeError)

				return m
			},
			expectedDTOs: nil,
			expectedErr:  sharedMocks.SomeError,
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
				m.EXPECT().ReadOne(sharedMocks.Entities[0].Id).Return(sharedMocks.Entities[0], nil)

				return m
			},
			entityId:    sharedMocks.Entities[0].Id,
			expectedDTO: sharedMocks.Entities[0],
			expectedErr: nil,
		},

		{
			name: "IdMissingInStorage error",
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().ReadOne(sharedMocks.Entities[0].Id).Return(shared.Entity{}, errs.IdMissingInStorage)

				return m
			},
			entityId:    sharedMocks.Entities[0].Id,
			expectedDTO: shared.Entity{},
			expectedErr: shared.NotFoundById,
		},

		{
			name: "general error",
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().ReadOne(sharedMocks.Entities[0].Id).Return(shared.Entity{}, sharedMocks.SomeError)

				return m
			},
			entityId:    sharedMocks.Entities[0].Id,
			expectedDTO: shared.Entity{},
			expectedErr: sharedMocks.SomeError,
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
			entityModel: sharedMocks.Entities[0],
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(sharedMocks.Entities[0].Id).Return(sharedMocks.AllEndpoints[:2], nil)
				m.EXPECT().Update(sharedMocks.Entities[0]).Return(sharedMocks.Entities[0], nil)

				return m
			},
			expectedDTO: sharedMocks.Entities[0],
			expectedErr: nil,
		},

		{
			name:        "IdMissingInStorage error",
			entityModel: sharedMocks.Entities[0],
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(sharedMocks.Entities[0].Id).Return(nil, errs.IdMissingInStorage)

				return m
			},
			expectedDTO: shared.Entity{},
			expectedErr: shared.NotFoundById,
		},

		{
			name:        "general error",
			entityModel: sharedMocks.Entities[0],
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(sharedMocks.Entities[0].Id).Return(nil, sharedMocks.SomeError)

				return m
			},
			expectedDTO: shared.Entity{},
			expectedErr: sharedMocks.SomeError,
		},

		{
			name:        "TryingToRemoveEndpointThatIsInUse error",
			entityModel: sharedMocks.Entities[0],
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(sharedMocks.Entities[0].Id).Return(sharedMocks.AllEndpoints[:4], nil)

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
			entityId: sharedMocks.Entities[0].Id,
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(sharedMocks.Entities[0].Id).Return(nil, nil)
				m.EXPECT().Delete(sharedMocks.Entities[0].Id).Return(nil)

				return m
			},
			expectedErr: nil,
		},

		{
			name:     "IdMissingInStorage error",
			entityId: sharedMocks.Entities[0].Id,
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(sharedMocks.Entities[0].Id).Return(nil, errs.IdMissingInStorage)

				return m
			},
			expectedErr: shared.NotFoundById,
		},

		{
			name:     "general error",
			entityId: sharedMocks.Entities[0].Id,
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(sharedMocks.Entities[0].Id).Return(nil, sharedMocks.SomeError)

				return m
			},
			expectedErr: sharedMocks.SomeError,
		},

		{
			name:     "TryingToRemoveEntityThatIsUse error",
			entityId: sharedMocks.Entities[0].Id,
			mockConstructor: func() Repo {
				m := entityMocks.NewMockRepo(t)
				m.EXPECT().FetchRelations(sharedMocks.Entities[0].Id).Return(sharedMocks.AllEndpoints[:3], nil)

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
