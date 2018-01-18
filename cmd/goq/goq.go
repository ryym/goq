package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"

	"github.com/ryym/goq/gen"
)

func main() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	pkgPath, err := filepath.Rel(filepath.Join(gopath, "src"), wd)
	if err != nil {
		panic(err)
	}
	opts := gen.Opts{
		Pkg:              pkgPath,
		OutFile:          "gql.go",
		TablesStructName: "Tables",
	}

	err = gen.GenerateTableHelpers(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
