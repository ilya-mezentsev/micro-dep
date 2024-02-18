package d2

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/frankenbeanies/uuid4"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/log"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"oss.terrastruct.com/util-go/go2"
)

type Client struct {
	dirname string
}

func New() Client {
	return Client{
		dirname: os.TempDir(), // can be configured if needed
	}
}

func (c Client) Draw(script string) (string, error) {
	ruler, err := textmeasure.NewRuler()
	if err != nil {
		return "", err
	}

	layoutResolver := func(engine string) (d2graph.LayoutGraph, error) {
		return func(ctx context.Context, g *d2graph.Graph) error {
			return d2dagrelayout.Layout(ctx, g, &d2dagrelayout.ConfigurableOpts{
				NodeSep: 60,
				EdgeSep: 20,
			})
		}, nil
	}
	renderOpts := &d2svg.RenderOpts{
		Pad:     go2.Pointer(int64(5)),
		ThemeID: &d2themescatalog.TerminalGrayscale.ID,
	}
	compileOpts := &d2lib.CompileOptions{
		LayoutResolver: layoutResolver,
		Ruler:          ruler,
	}

	diagram, _, err := d2lib.Compile(
		log.With(
			context.Background(),
			slog.Make(sloghuman.Sink(io.Discard)),
		),
		script,
		compileOpts,
		renderOpts,
	)
	if err != nil {
		return "", err
	}

	out, err := d2svg.Render(diagram, renderOpts)
	if err != nil {
		return "", err
	}

	path := filepath.Join(c.dirname, uuid4.New().String()+".svg")
	err = os.WriteFile(path, out, 0600)

	return path, err
}
