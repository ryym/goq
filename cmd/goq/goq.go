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
	isTestPkg := flag.Bool("testpkg", false, "Generate from code in test package")
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
		PkgPath:          pkgPath,
		OutFile:          outFile,
		TablesStructName: "Tables",
		IsTestPkg:        *isTestPkg,
	}

	err = gen.GenerateTableHelpers(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
