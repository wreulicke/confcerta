package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/wreulicke/confcerta/backend"
	"github.com/wreulicke/confcerta/backend/internal"
)

type Backend struct {
	client   s3iface.S3API
	bucket   string
	key      string
	Splitter backend.Splitter
}

func New(s3 s3iface.S3API, bucket, path string) backend.Backend {
	return &Backend{
		client:   s3,
		bucket:   bucket,
		key:      path,
		Splitter: backend.NewSplitter("/"),
	}
}

func (b *Backend) Load(ctx context.Context, cfg *backend.Config) (map[string]interface{}, error) {
	r := map[string]interface{}{}
	o, err := b.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &b.bucket,
		Key:    &b.key,
	})
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == s3.ErrCodeNoSuchKey {
			return r, nil
		}
		return nil, err
	} else if err != nil {
		return nil, err
	}
	defer o.Body.Close()

	d, err := internal.NewDecoder(b.key, o.Body)
	if err != nil {
		return nil, err
	}
	return r, d(r)
}
