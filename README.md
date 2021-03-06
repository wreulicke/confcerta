## Confcerta 

[![PkgGoDev](https://pkg.go.dev/badge/wreulicke/confcerta)](https://pkg.go.dev/wreulicke/confcerta)

Confcerta is a library that load configuration from multiple backends, inspired by [confita](https://github.com/heetch/confita).

## Usage

```go
import (
    "github.com/wreulicke/confcerta"
)

type Config struct {
	A      string 
	// You can specify aliases using tag 
	B      string `config:"Foo"`
	C      string `config:"Bar"`
}

func main() {
    // json
    // {"A": "A value", "Foo": "B value", "Bar": "C value"}
    l := confcerta.New(file.New("testdata/simple.json"))
    c := Config{}
	l.Unmarshal(context.Background(), &c)
    // c.A == "A value"
    // c.B == "B value"
    // c.C == "C value"
}
```

## Supported backends

- Environment variables
- Several format of files
  - JSON files
  - Yaml files
  - Toml files
  - HCL files
- Command line flags
- Amazon Systems Manager Parameter Store
* Amazon S3
* http.FileSystem
* io.Reader

## Install

```sh
go get -u github.com/wreulicke/confcerta
```
