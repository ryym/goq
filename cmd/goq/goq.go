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

	var outFile string
	args := flag.Args()
	switch len(args) {
	case 0:
		outFile = "gql.go"
	case 1:
		if strings.HasSuffix(args[0], ".go") {
			if strings.ContainsRune(args[0], filepath.Separator) {
				panic("output file must be in the same directory")
			}
			outFile = args[0]
		} else {
			panic("Invalid output file name")
		}
	default:
		panic("Invalid arguments")
	}

	opts := gen.Opts{
		Pkg:              pkgPath,
		OutFile:          outFile,
		TablesStructName: "Tables",
		ImportTests:      *withTests,
	}

	err = gen.GenerateTableHelpers(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
