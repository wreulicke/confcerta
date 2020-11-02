package internal

import (
	"encoding/json"
	"errors"
	"io"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type Decoder func(to interface{}) error

func NewDecoderFromFile(path string, reader io.Reader) (Decoder, error) {
	return NewDecoderFromExtension(filepath.Ext(path), reader)
}

func NewDecoder(p string, reader io.Reader) (Decoder, error) {
	return NewDecoderFromExtension(path.Ext(p), reader)
}

func NewDecoderFromExtension(ext string, reader io.Reader) (Decoder, error) {
	switch ext {
	case ".json":
		return json.NewDecoder(reader).Decode, nil
	case ".yml":
		fallthrough
	case ".yaml":
		return yaml.NewDecoder(reader).Decode, nil
	case ".toml":
		return func(to interface{}) error {
			_, err := toml.DecodeReader(reader, to)
			return err
		}, nil
	}
	return nil, errors.New("Unsupported format")
}
