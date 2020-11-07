package confcerta

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/wreulicke/confcerta/backend/env"
	"github.com/wreulicke/confcerta/backend/file"
	"github.com/wreulicke/confcerta/backend/flags"
)

type testConfig struct {
	A      string `config:"alice"`
	E      string `config:"eugeo"`
	K      string `config:"kirito"`
	Nested struct {
		Key   string
		Value string `config:"value"`
	}
}

func TestLoad_JSON(t *testing.T) {
	l := New(file.New("testdata/simple.json"))
	c := testConfig{}
	err := l.Unmarshal(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.A != "Alice Zuberg" {
		t.Error("c.A is not alice")
	}
	if c.E != "Eugeo" {
		t.Error("c.E is not eugeo")
	}
	if c.K != "Kazuto Kirigaya" {
		t.Error("c.K is not kirito")
	}
	if c.Nested.Key != "Nested Key" {
		t.Error("c.Nested.Key is not expected condition")
	}
	if c.Nested.Value != "Nested Value" {
		t.Error("c.Nested.Value is not expected condition")
	}
}

func TestLoad_EnvBackned(t *testing.T) {
	os.Setenv("ALICE", "Alice Zuberg")
	os.Setenv("Eugeo", "Eugeo")
	os.Setenv("Kirito", "Kazuto Kirigaya")
	os.Setenv("NESTED_KEY", "Nested Key")
	os.Setenv("NESTED_VALUE", "Nested Value")
	l := New(env.New())
	c := testConfig{}
	err := l.Unmarshal(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.A != "Alice Zuberg" {
		t.Error("c.A is not alice")
	}
	if c.E != "Eugeo" {
		t.Error("c.E is not eugeo")
	}
	if c.K != "Kazuto Kirigaya" {
		t.Error("c.K is not kirito")
	}
	if c.Nested.Key != "Nested Key" {
		t.Error("c.Nested.Key is not expected condition")
	}
	if c.Nested.Value != "Nested Value" {
		t.Error("c.Nested.Value is not expected condition")
	}
}

func TestLoad_FlagsBackned(t *testing.T) {
	os.Args = []string{"",
		"-alice", "Alice Zuberg",
		"-eugeo", "Eugeo",
		"-kirito", "Kazuto Kirigaya",
		"-nested-key", "Nested Key",
		"-nested-value", "Nested Value",
	}
	l := New(flags.New())
	c := testConfig{}
	err := l.Unmarshal(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	if c.A != "Alice Zuberg" {
		t.Error("c.A is not alice")
	}
	if c.E != "Eugeo" {
		t.Error("c.E is not eugeo")
	}
	if c.K != "Kazuto Kirigaya" {
		t.Error("c.K is not kirito")
	}
	if c.Nested.Key != "Nested Key" {
		t.Error("c.Nested.Key is not expected condition")
	}
	if c.Nested.Value != "Nested Value" {
		t.Error("c.Nested.Value is not expected condition")
	}
}

func TestLoad_HCL(t *testing.T) {
	l := New(file.New("testdata/simple.hcl"))
	c := struct {
		Character []struct {
			Name string `config:"name"`
		} `config:"character"`
	}{}
	err := l.Unmarshal(context.Background(), &c)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(c)
	if c.Character[0].Name != "Alice Zuberg" {
		t.Error("c.A is not alice")
	}
	if c.Character[1].Name != "Eugeo" {
		t.Error("c.E is not eugeo")
	}
	if c.Character[2].Name != "Kazuto Kirigaya" {
		t.Error("c.K is not kirito")
	}
}
