package file

import (
	"context"
	"os"
	"strings"

	"github.com/wreulicke/confcerta/backend"
	"github.com/wreulicke/confcerta/backend/internal"
)

type Backend struct {
	path     string
	optional bool
}

func expandPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return os.ExpandEnv(strings.Replace(path, "~", home, 1))
}

func New(path string) backend.Backend {
	return &Backend{
		path: expandPath(path),
	}
}

func NewOptional(path string) backend.Backend {
	return &Backend{
		path:     expandPath(path),
		optional: true,
	}
}

func (b *Backend) Load(ctx context.Context, cfg *backend.Config) (map[string]interface{}, error) {
	r := map[string]interface{}{}
	f, err := os.Open(b.path)
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
