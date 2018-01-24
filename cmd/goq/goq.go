package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/ryym/goq/gen"
)

func main() {
	withTests := flag.Bool("tests", false, "Include test files for parsing")
	flag.Parse()

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
	args := flag.Args()
	switch len(args) {
	case 0:
		outPath = "gql.go"
	case 1:
		if strings.HasSuffix(args[0], ".go") {
			outPath = args[0]
		} else {
			panic("Invalid output file name")
		}
	default:
		panic("Invalid arguments")
	}

	opts := gen.Opts{
		Pkg:              pkgPath,
		OutPath:          outPath,
		TablesStructName: "Tables",
		ImportTests:      *withTests,
	}

	err = gen.GenerateTableHelpers(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
