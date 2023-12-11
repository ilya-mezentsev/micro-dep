package relation

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	relationMocks "github.com/ilya-mezentsev/micro-dep/store/internal/services/relation/mocks"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
	sharedMocks "github.com/ilya-mezentsev/micro-dep/store/internal/services/shared/mocks"
)

func TestServiceImpl_Create(t *testing.T) {
	tests := []struct {
		name            string
		model           shared.Relation
		mockConstructor func() Repo
		expected        error
	}{
		{
			name:  "ok",
			model: sharedMocks.Relations[0],
			mockConstructor: func() Repo {
				m := relationMocks.NewMockRepo(t)
				m.EXPECT().PartsExist(sharedMocks.Relations[0]).Return(true, true, nil)
				m.EXPECT().Create(mock.Anything).Return(sharedMocks.Relations[0], nil)

				return m
			},
			expected: nil,
		},

		{
			name:  "IdMissingInStorage error",
			model: sharedMocks.Relations[0],
			mockConstructor: func() Repo {
				m := relationMocks.NewMockRepo(t)
				m.EXPECT().PartsExist(sharedMocks.Relations[0]).Return(false, false, errs.IdMissingInStorage)

				return m
			},
			expected: shared.NotFoundById,
		},

		{
			name:  "some error",
			model: sharedMocks.Relations[0],
			mockConstructor: func() Repo {
				m := relationMocks.NewMockRepo(t)
				m.EXPECT().PartsExist(sharedMocks.Relations[0]).Return(false, false, sharedMocks.SomeError)

				return m
			},
			expected: sharedMocks.SomeError,
		},

		{
			name:  "FROM entity is missed error",
			model: sharedMocks.Relations[0],
			mockConstructor: func() Repo {
				m := relationMocks.NewMockRepo(t)
				m.EXPECT().PartsExist(sharedMocks.Relations[0]).Return(false, false, nil)

				return m
			},
			expected: TryingToCreateRelationFromMissedEntity,
		},

		{
			name:  "TO endpoint is missed error",
			model: sharedMocks.Relations[0],
			mockConstructor: func() Repo {
				m := relationMocks.NewMockRepo(t)
				m.EXPECT().PartsExist(sharedMocks.Relations[0]).Return(true, false, nil)

				return m
			},
			expected: TryingToCreateRelationToMissedEndpoint,
		},

		{
			name:  "creation error",
			model: sharedMocks.Relations[0],
			mockConstructor: func() Repo {
				m := relationMocks.NewMockRepo(t)
				m.EXPECT().PartsExist(sharedMocks.Relations[0]).Return(true, true, nil)
				m.EXPECT().Create(mock.Anything).Return(shared.Relation{}, sharedMocks.SomeError)

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

func TestServiceImpl_ReadAll(t *testing.T) {
	tests := []struct {
		name            string
		mockConstructor func() Repo
		expectedModels  []shared.Relation
		expectedError   error
	}{
		{
			name: "ok",
			mockConstructor: func() Repo {
				m := relationMocks.NewMockRepo(t)
				m.EXPECT().ReadAll().Return(sharedMocks.Relations[:1], nil)

				return m
			},
			expectedModels: sharedMocks.Relations[:1],
			expectedError:  nil,
		},

		{
			name: "some error",
			mockConstructor: func() Repo {
				m := relationMocks.NewMockRepo(t)
				m.EXPECT().ReadAll().Return(nil, sharedMocks.SomeError)

				return m
			},
			expectedModels: nil,
			expectedError:  sharedMocks.SomeError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServiceImpl(tt.mockConstructor())
			allModels, err := s.ReadAll()

			require.Equal(t, tt.expectedModels, allModels)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestServiceImpl_ReadOne(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()

	_, _ = NewServiceImpl(relationMocks.NewMockRepo(t)).ReadOne(sharedMocks.Relations[0].Id)
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
			modelId: sharedMocks.Relations[0].Id,
			mockConstructor: func() Repo {
				m := relationMocks.NewMockRepo(t)
				m.EXPECT().Delete(sharedMocks.Relations[0].Id).Return(nil)

				return m
			},
			expected: nil,
		},

		{
			name:    "some error",
			modelId: sharedMocks.Relations[0].Id,
			mockConstructor: func() Repo {
				m := relationMocks.NewMockRepo(t)
				m.EXPECT().Delete(sharedMocks.Relations[0].Id).Return(sharedMocks.SomeError)

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
