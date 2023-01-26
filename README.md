# Go CI

![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/sagikazarmark/goci/ci.yaml?style=flat-square)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.19-61CFDD.svg?style=flat-square)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/mod/github.com/sagikazarmark/goci)
[![built with nix](https://img.shields.io/badge/builtwith-nix-7d81f7?style=flat-square)](https://builtwithnix.org)

**High-level CI config DSL written in Go based on [Dagger](https://dagger.io/).**

**⚠️ This tool is still under heavy development! Things may change. ⚠️**

## Features

- [x] Go test
- [x] GolangCI Lint
- [x] CI detection
- [ ] CodeCov upload
- [ ] Build pipelines
  - [ ] Build matrix
  - [ ] Pipelines
  - [ ] Step dependencies

## Goals

- Create a high-level interface for building a CI based on Dagger
- Hide low-level (Dagger) details as much as possible

## Usage

Install the library:

```shell
go get github.com/sagikazarmark/goci
```

Create CLI tool:

```go
package main

func main() {
	client, err := dagger.Connect(ctx)
	if err != nil {
		return panic(err)
	}
	defer client.Close()

	c := golang.Test(client)

	output, err := container.Stdout(ctx)
	if err != nil {
		return panic(err)
	}

	fmt.Print(output)
}
```


## License

The project is licensed under the [MIT License](LICENSE).
