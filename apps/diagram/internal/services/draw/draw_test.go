package draw

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	drawMocks "github.com/ilya-mezentsev/micro-dep/diagram/internal/services/draw/mocks"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/shared/types"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

var rdp = types.RelationsDiagramData{
	Entities: []models.Entity{
		{
			Id:   "foo",
			Name: "foo",
			Endpoints: []models.Endpoint{
				{
					Id:       "foo-e-1",
					EntityId: "foo",
					Kind:     "foo e 1 endpoint",
					Address:  "/api/foo/e/1",
				},
				{
					Id:       "foo-e-2",
					EntityId: "foo",
					Kind:     "foo e 2 endpoint",
					Address:  "/api/foo/e/2",
				},
			},
		},

		{
			Id:   "bar",
			Name: "bar",
			Endpoints: []models.Endpoint{
				{
					Id:       "bar-e-1",
					EntityId: "bar",
					Kind:     "bar e 1 endpoint",
					Address:  "/api/bar/e/1",
				},
				{
					Id:       "bar-e-2",
					EntityId: "bar",
					Kind:     "bar e 2 endpoint",
					Address:  "/api/bar/e/2",
				},
			},
		},

		{
			Id:   "baz",
			Name: "baz",
			Endpoints: []models.Endpoint{
				{
					Id:       "baz-e-1",
					EntityId: "baz",
					Kind:     "baz e 1 endpoint",
					Address:  "/api/baz/e/1",
				},
				{
					Id:       "baz-e-2",
					EntityId: "baz",
					Kind:     "baz e 2 endpoint",
					Address:  "/api/baz/e/2",
				},
			},
		},

		{
			Id:   "xyz",
			Name: "xyz",
		},
	},
	Relations: []models.Relation{
		{
			FromEntityId: "foo",
			ToEndpointId: "bar-e-1",
		},
		{
			FromEntityId: "foo",
			ToEndpointId: "baz-e-2",
		},
		{
			FromEntityId: "bar",
			ToEndpointId: "baz-e-1",
		},
		{
			FromEntityId: "xyz",
			ToEndpointId: "bar-e-2",
		},
	},
}

const expectedResult = `

	foo: {
		
			foo-e-1: /api/foo/e/1
		
			foo-e-2: /api/foo/e/2
		
	}

	bar: {
		
			bar-e-1: /api/bar/e/1
		
			bar-e-2: /api/bar/e/2
		
	}

	baz: {
		
			baz-e-1: /api/baz/e/1
		
			baz-e-2: /api/baz/e/2
		
	}

	xyz: {
		
	}




	foo -> bar.bar-e-1 : bar e 1 endpoint

	foo -> baz.baz-e-2 : baz e 2 endpoint

	bar -> baz.baz-e-1 : baz e 1 endpoint

	xyz -> bar.bar-e-2 : bar e 2 endpoint

`

func TestService_DrawDiagramOk(t *testing.T) {
	m := drawMocks.NewMockD2Client(t)
	m.EXPECT().Draw(expectedResult).Return("foo", nil)

	s := New(m)

	res, err := s.DrawDiagram(rdp)

	require.Equal(t, "foo", res)
	require.Nil(t, err)
}

func TestService_DrawDiagramErr(t *testing.T) {
	someErr := errors.New("some-err")

	m := drawMocks.NewMockD2Client(t)
	m.EXPECT().Draw(expectedResult).Return("foo", someErr)

	s := New(m)

	res, err := s.DrawDiagram(rdp)

	require.Equal(t, "foo", res)
	require.Equal(t, someErr, err)
}
