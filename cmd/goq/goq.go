package main

import (
	"errors"
	"flag"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/ryym/goq/gen"
)

func main() {
	IsTestFile := flag.Bool("test", false, "Generate from code in test file")
	flag.Parse()

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	wd, err := os.Getwd()
	if err != nil {
		exitErr(err)
	}

	pkgPath, err := filepath.Rel(filepath.Join(gopath, "src"), wd)
	if err != nil {
		exitErr(err)
	}

	var outFile string
	args := flag.Args()
	switch len(args) {
	case 0:
		outFile = "gql.go"
	case 1:
		if strings.ContainsRune(args[0], filepath.Separator) {
			exitErr(errors.New("output file must be in the same directory"))
		}
		outFile = args[0]
	default:
		exitErr(errors.New("cannot specify more than 1 file"))
	}

	opts := gen.Opts{
		PkgPath:          pkgPath,
		OutFile:          outFile,
		TablesStructName: "Tables",
		IsTestFile:       *IsTestFile,
	}

	err = gen.GenerateTableHelpers(opts)
	if err != nil {
		exitErr(err)
	}
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
