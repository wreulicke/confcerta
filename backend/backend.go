package backend

import (
	"context"
	"strings"
)

type Splitter func(string) []string

type Config struct {
	Fields []*Field
}

type Field struct {
	Name string
}

type Backend interface {
	Load(ctx context.Context, cfg *Config) (map[string]interface{}, error)
}

func NewSplitter(sep string) Splitter {
	return func(s string) []string {
		return strings.Split(s, sep)
	}
}
