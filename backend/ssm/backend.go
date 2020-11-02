package ssm

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/wreulicke/confcerta/backend"
)

type Backend struct {
	client   ssmiface.SSMAPI
	path     string
	Splitter backend.Splitter
}

func New(ssm ssmiface.SSMAPI, path string) backend.Backend {
	return &Backend{
		client:   ssm,
		path:     path,
		Splitter: backend.NewSplitter("/"),
	}
}

func (b *Backend) Load(ctx context.Context, cfg *backend.Config) (map[string]interface{}, error) {
	r := map[string]interface{}{}
	recursive := true
	decryption := true
	var max int64 = 10
	input := &ssm.GetParametersByPathInput{
		Path:           &b.path,
		Recursive:      &recursive,
		WithDecryption: &decryption,
		MaxResults:     &max,
	}

	for {
		res, err := b.client.GetParametersByPathWithContext(ctx, input)
		if err != nil {
			return nil, err
		}
		for _, p := range res.Parameters {
			if p.Name == nil || p.Value == nil {
				continue
			}
			path := strings.TrimPrefix(*p.Name, b.path)
			pathElems := b.Splitter(path)
			target := r
			for index, name := range pathElems {
				if index < len(pathElems)-1 {
					if _, ok := target[name]; !ok {
						v := map[string]interface{}{}
						target[name] = v
						target = v
					} else if v, ok := target[name].(map[string]interface{}); ok {
						target = v
					} else {
						err = fmt.Errorf("cannot allocate field. path: %s", *p.Name)
						return nil, err
					}
				} else {
					target[name] = *p.Value
				}
			}
		}
		if res.NextToken == nil {
			break
		}
		input.NextToken = res.NextToken
	}
	return r, nil
}

func newBool(b bool) *bool {
	return &b
}

func newInt64(i int64) *int64 {
	return &i
}
