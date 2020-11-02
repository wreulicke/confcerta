package confcerta

import (
	"context"
	"errors"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/wreulicke/confcerta/backend"
)

type loader struct {
	backends []backend.Backend
}

type Unmarshaler interface {
	Unmarshal(ctx context.Context, to interface{}) error
}

func New(backends ...backend.Backend) Unmarshaler {
	return &loader{
		backends: backends,
	}
}

func concatPath(parent string, nested string) string {
	if parent == "" {
		return nested
	}
	return parent + "-" + nested
}

func (l *loader) parseStruct(ctx context.Context, path string, ref reflect.Value) (*backend.Config, error) {
	c := &backend.Config{}

	t := ref.Type()
	numFields := ref.NumField()
	for i := 0; i < numFields; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		field := t.Field(i)
		value := ref.Field(i)
		typ := value.Type()

		name := concatPath(path, field.Name)
		tag := field.Tag.Get("config")
		if tag != "" && tag != "-" {
			name = concatPath(path, tag)
		}

		switch typ.Kind() {
		case reflect.Struct:
			p, err := l.parseStruct(ctx, name, value)
			if err != nil {
				return nil, err
			}
			c.Fields = append(c.Fields, p.Fields...)
			continue
		case reflect.Ptr:
			if typ.Elem().Kind() == reflect.Struct && !value.IsNil() {
				p, err := l.parseStruct(ctx, name, value.Elem())
				if err != nil {
					return nil, err
				}
				c.Fields = append(c.Fields, p.Fields...)
				continue
			}
		}

		f := &backend.Field{
			Name: name,
		}
		c.Fields = append(c.Fields, f)
	}

	return c, nil
}

func (l *loader) Unmarshal(ctx context.Context, to interface{}) error {
	ref := reflect.ValueOf(to)
	if !ref.IsValid() || ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Struct {
		return errors.New("target must be a pointer to struct")
	}
	ref = ref.Elem()
	config, err := l.parseStruct(ctx, "", ref)
	if err != nil {
		return err
	}

	mapstructureConfig := &mapstructure.DecoderConfig{
		TagName: "config",
		Result:  to,
	}
	d, err := mapstructure.NewDecoder(mapstructureConfig)
	if err != nil {
		return err
	}

	r := map[string]interface{}{}
	for _, e := range l.backends {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		m, err := e.Load(ctx, config)
		if err != nil {
			return err
		}
		for key, value := range m {
			r[key] = value
		}
	}

	return d.Decode(r)
}
