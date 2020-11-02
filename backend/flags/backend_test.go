package flags

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/wreulicke/confcerta/backend"
)

func Test(t *testing.T) {
	b := Backend{
		Flags:    &flag.FlagSet{},
		Splitter: backend.NewSplitter("-"),
	}
	os.Args = []string{"", "-simple", "simple-value", "-nested-field", "nested-value"}
	r, err := b.Load(context.Background(), &backend.Config{
		Fields: []*backend.Field{
			{
				Name: "simple",
			},
			{
				Name: "nested-field",
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
	if r["simple"] != "simple-value" {
		t.Error("path `simple` is not able to parse properly")
	}
	if r["nested"].(map[string]interface{})["field"] != "nested-value" {
		t.Error("path `nested-field` is not able to parse properly")
	}
}
