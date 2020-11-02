package env

import (
	"context"
	"os"
	"strings"

	"github.com/wreulicke/confcerta/backend"
)

type Backend struct {
	Splitter backend.Splitter
}

func New() backend.Backend {
	return &Backend{
		Splitter: backend.NewSplitter("_"),
	}
}

func (b *Backend) Load(ctx context.Context, cfg *backend.Config) (map[string]interface{}, error) {
	r := map[string]interface{}{}
	for _, kv := range os.Environ() {
		arr := strings.SplitN(kv, "=", 2)
		key := arr[0]
		value := arr[1]
		p := b.Splitter(strings.ToLower(key))
		target := r
		for index, name := range p {
			if index < len(p)-1 {
				if _, ok := target[name]; !ok {
					v := map[string]interface{}{}
					target[name] = v
					target = v
				} else if v, ok := target[name].(map[string]interface{}); ok {
					target = v
				} else {
					// ignore to fail on nested key
					continue
				}
			} else {
				target[name] = value
			}
		}
	}
	return r, nil
}
