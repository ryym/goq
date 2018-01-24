package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

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

	var outPath string
	if len(os.Args) == 2 && strings.HasSuffix(os.Args[1], ".go") {
		outPath = os.Args[1]
	} else {
		outPath = "gql.go"
	}
	opts := gen.Opts{
		Pkg:              pkgPath,
		OutPath:          outPath,
		TablesStructName: "Tables",
	}

	err = gen.GenerateTableHelpers(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
