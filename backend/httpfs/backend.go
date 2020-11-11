package httpfs

import (
	"context"
	"net/http"
	"os"

	"github.com/wreulicke/confcerta/backend"
	"github.com/wreulicke/confcerta/backend/internal"
)

type Backend struct {
	http.FileSystem
	path     string
	optional bool
}

func New(fs http.FileSystem, path string) backend.Backend {
	return &Backend{
		FileSystem: fs,
		path:       expandPath(path),
	}
}

func NewOptional(fs http.FileSystem, path string) backend.Backend {
	return &Backend{
		FileSystem: fs,
		path:       expandPath(path),
		optional:   true,
	}
}

func expandPath(path string) string {
	return os.ExpandEnv(path)
}

func (b *Backend) Load(ctx context.Context, cfg *backend.Config) (map[string]interface{}, error) {
	r := map[string]interface{}{}
	f, err := b.Open(b.path)
	if err != nil {
		if os.IsNotExist(err) && b.optional {
			return r, nil
		}
		return nil, err
	}

	defer f.Close()
	d, err := internal.NewDecoderFromFile(b.path, f)
	if err != nil {
		return nil, err
	}
	return r, d(&r)
}
