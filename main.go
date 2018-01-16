package main

import (
	"fmt"
	"os"
	"path"

	"github.com/ryym/goq/gen"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	opts := gen.Opts{
		DestDir:          path.Join(wd, "gql"),
		ModelsDir:        wd,
		TablesStructName: "Tables",
	}

	err = gen.GenerateTableHelpers(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
