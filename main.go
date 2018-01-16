package main

import (
	"fmt"
	"os"

	"github.com/ryym/goq/gen"
)

func main() {
	opts := gen.Opts{
		OutFile:          "gql.go",
		TablesStructName: "Tables",
	}

	err := gen.GenerateTableHelpers(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
