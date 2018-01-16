package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/ryym/goq/util"
)

type Opts struct {
	DestDir          string
	ModelsDir        string
	TablesStructName string
}

type helper struct {
	Name      string
	TableName string
	ModelName string
	Fields    []*field
}

func (h *helper) JoinFields(alias string) string {
	cols := make([]string, len(h.Fields))
	for i, f := range h.Fields {
		cols[i] = fmt.Sprintf("%s.%s", alias, f.Name)
	}
	return strings.Join(cols, ", ")
}

type field struct {
	Name   string
	Column string
}

func GenerateTableHelpers(opts Opts) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, opts.ModelsDir, exceptTests, 0)
	if err != nil {
		panic(err)
	}

	var pkgName string
	for name, _ := range pkgs {
		pkgName = name
		break
	}

	structs := findAllStructs(pkgs)

	tableList, ok := structs[opts.TablesStructName]
	if !ok {
		return errors.New("Tables struct not found")
	}

	helpers := make([]*helper, len(tableList.Fields.List))
	for i, table := range tableList.Fields.List {
		if len(table.Names) == 0 {
			continue // TODO: Support embedded structs.
		}

		tableName := table.Names[0].Name
		modelName := table.Type.(*ast.Ident).Name
		if _, ok := structs[modelName]; !ok {
			return fmt.Errorf("model %s not found", modelName)
		}

		helpers[i] = &helper{
			Name:      util.ColToFld(tableName),
			TableName: tableName,
			ModelName: modelName,
			Fields:    listColumnFields(structs[modelName]),
		}
	}

	// dirStat, err := os.Stat("gql")
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		err = os.Mkdir("gql", 0766)
	// 	}
	// 	if err != nil {
	// 		return errors.Wrap("failed to create sub package", err)
	// 	}
	// } else if !dirStat.IsDir() {
	// 	return errors.New("gql exists but is not a directory")
	// }

	outPath := filepath.Join("gql.go")
	if _, err := os.Stat(outPath); err == nil {
		err = os.Remove(outPath)
		if err != nil {
			return fmt.Errorf("failed to remove %s", outPath)
		}
	}

	file, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create %s", outPath)
	}
	defer file.Close()

	err = writeTemplate(file, pkgName, helpers)
	if err != nil {
		return err
	}

	return nil
}

func exceptTests(f os.FileInfo) bool {
	return !strings.HasSuffix(f.Name(), "_test.go")
}

func findAllStructs(pkgs map[string]*ast.Package) map[string]*ast.StructType {
	structs := map[string]*ast.StructType{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, dcl := range file.Decls {
				if d, ok := dcl.(*ast.GenDecl); ok && d.Tok == token.TYPE {
					for _, sp := range d.Specs {
						tsp := sp.(*ast.TypeSpec)
						if ast.IsExported(tsp.Name.Name) {
							if st, ok := tsp.Type.(*ast.StructType); ok {
								structs[tsp.Name.Name] = st
							}
						}
					}
				}
			}
		}
	}
	return structs
}

func listColumnFields(model *ast.StructType) []*field {
	var fields []*field
	for _, fld := range model.Fields.List {
		fldName := fld.Names[0].Name
		if ast.IsExported(fldName) {
			fields = append(fields, &field{
				Name:   fldName,
				Column: util.FldToCol(fldName),
			})
		}
	}
	return fields
}
