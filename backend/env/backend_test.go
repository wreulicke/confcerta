package env

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/wreulicke/confcerta/backend"
)

func TestBackend_Load(t *testing.T) {
	os.Setenv("SIMPLE", "test")
	os.Setenv("NESTED_KEY", "nested-value")
	type args struct {
		cfg *backend.Config
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				cfg: &backend.Config{
					Fields: []*backend.Field{
						{Name: "simple"},
					},
				},
			},
			want: map[string]interface{}{
				"simple": "test",
			},
		},
		{
			name: "nested",
			args: args{
				cfg: &backend.Config{
					Fields: []*backend.Field{
						{Name: "nested-key"},
					},
				},
			},
			want: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": "nested-value",
				},
			},
		},
	}
	b := &Backend{
		Splitter: backend.NewSplitter("_"),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := b.Load(context.Background(), tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Backend.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for k, v := range tt.want {
				if !reflect.DeepEqual(got[k], v) {
					t.Errorf("result[%s] = %v, want %v", k, got[k], v)
				}
			}
		})
	}
}
