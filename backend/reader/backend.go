package reader

import (
	"context"
	"io"

	"github.com/wreulicke/confcerta/backend"
	"github.com/wreulicke/confcerta/backend/internal"
)

type Backend struct {
	reader io.Reader
	ext    string
}

func New(reader io.Reader, ext string) backend.Backend {
	return &Backend{
		reader: reader,
		ext:    ext,
	}
}

func (b *Backend) Load(ctx context.Context, cfg *backend.Config) (map[string]interface{}, error) {
	r := map[string]interface{}{}
	d, err := internal.NewDecoderFromExtension(b.ext, b.reader)
	if err != nil {
		return nil, err
	}
	return r, d(&r)
}
