package diagram

import (
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	diagramMocks "github.com/ilya-mezentsev/micro-dep/diagram/internal/services/diagram/mocks"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/shared/types"
	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

var (
	someAccountId   models.Id = "foo-bar-baz"
	someDiagramPath           = "/tmp/xyz.svg"
	someError                 = errors.New("some-error")
)

func TestService_Draw(t *testing.T) {
	tests := []struct {
		name                string
		accountId           models.Id
		mocksConstructor    func() (EntitiesFetcher, RelationsFetcher, DrawService)
		expectedDiagramPath string
		expectedError       error
	}{
		{
			name:      "ok",
			accountId: someAccountId,
			mocksConstructor: func() (EntitiesFetcher, RelationsFetcher, DrawService) {
				ef := diagramMocks.NewMockEntitiesFetcher(t)
				ef.EXPECT().Fetch(someAccountId).Return(nil, nil)

				rf := diagramMocks.NewMockRelationsFetcher(t)
				rf.EXPECT().Fetch(someAccountId).Return(nil, nil)

				ds := diagramMocks.NewMockDrawService(t)
				ds.EXPECT().DrawDiagram(types.RelationsDiagramData{
					Entities:  nil,
					Relations: nil,
				}).Return(someDiagramPath, nil)

				return ef, rf, ds
			},
			expectedDiagramPath: someDiagramPath,
			expectedError:       nil,
		},

		{
			name:      "error from entities fetcher",
			accountId: someAccountId,
			mocksConstructor: func() (EntitiesFetcher, RelationsFetcher, DrawService) {
				ef := diagramMocks.NewMockEntitiesFetcher(t)
				ef.EXPECT().Fetch(someAccountId).Return(nil, someError)

				rf := diagramMocks.NewMockRelationsFetcher(t)
				rf.EXPECT().Fetch(someAccountId).Return(nil, nil)

				return ef, rf, nil
			},
			expectedError: errors.Join(someError, nil),
		},

		{
			name:      "error from relations fetcher",
			accountId: someAccountId,
			mocksConstructor: func() (EntitiesFetcher, RelationsFetcher, DrawService) {
				ef := diagramMocks.NewMockEntitiesFetcher(t)
				ef.EXPECT().Fetch(someAccountId).Return(nil, nil)

				rf := diagramMocks.NewMockRelationsFetcher(t)
				rf.EXPECT().Fetch(someAccountId).Return(nil, someError)

				return ef, rf, nil
			},
			expectedError: errors.Join(nil, someError),
		},

		{
			name:      "error from entities AND relations fetcher",
			accountId: someAccountId,
			mocksConstructor: func() (EntitiesFetcher, RelationsFetcher, DrawService) {
				ef := diagramMocks.NewMockEntitiesFetcher(t)
				ef.EXPECT().Fetch(someAccountId).Return(nil, someError)

				rf := diagramMocks.NewMockRelationsFetcher(t)
				rf.EXPECT().Fetch(someAccountId).Return(nil, someError)

				return ef, rf, nil
			},
			expectedError: errors.Join(someError, someError),
		},

		{
			name:      "error from draw service",
			accountId: someAccountId,
			mocksConstructor: func() (EntitiesFetcher, RelationsFetcher, DrawService) {
				ef := diagramMocks.NewMockEntitiesFetcher(t)
				ef.EXPECT().Fetch(someAccountId).Return(nil, nil)

				rf := diagramMocks.NewMockRelationsFetcher(t)
				rf.EXPECT().Fetch(someAccountId).Return(nil, nil)

				ds := diagramMocks.NewMockDrawService(t)
				ds.EXPECT().DrawDiagram(types.RelationsDiagramData{
					Entities:  nil,
					Relations: nil,
				}).Return("", someError)

				return ef, rf, ds
			},
			expectedDiagramPath: "",
			expectedError:       errs.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ef, rf, ds := tt.mocksConstructor()
			s := New(ef, rf, ds, slog.New(slog.NewTextHandler(io.Discard, nil)))

			diagramPath, err := s.Draw(tt.accountId)

			require.Equal(t, tt.expectedDiagramPath, diagramPath)
			require.Equal(t, tt.expectedError, err)
		})
	}
}
