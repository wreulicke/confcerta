package flags

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/wreulicke/confcerta/backend"
)

type Backend struct {
	Flags    *flag.FlagSet
	Splitter backend.Splitter
}

func New() backend.Backend {
	return &Backend{
		Flags:    flag.CommandLine,
		Splitter: backend.NewSplitter("-"),
	}
}

func (b *Backend) Load(ctx context.Context, cfg *backend.Config) (map[string]interface{}, error) {
	for _, e := range cfg.Fields {
		b.Flags.String(strings.ToLower(e.Name), "", e.Name)
	}
	err := b.Flags.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	r := map[string]interface{}{}
	b.Flags.VisitAll(func(f *flag.Flag) {
		p := b.Splitter(f.Name)
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
					err = fmt.Errorf("cannot allocate field. path: %s", f.Name)
					return
				}
			} else {
				target[name] = f.Value.String()
			}
		}
	})
	return r, err
}
