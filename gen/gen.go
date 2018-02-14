package gen

import (
	"fmt"
	"go/types"
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/ryym/goq/util"
	"golang.org/x/tools/go/loader"
)

type Opts struct {
	PkgPath          string
	OutFile          string
	TablesStructName string
	IsTestFile       bool
}

type helper struct {
	Name         string
	TableName    string
	ModelPkgName string
	ModelName    string
	Fields       []*field
}

func (h *helper) JoinFields(alias string) string {
	cols := make([]string, len(h.Fields))
	for i, f := range h.Fields {
		cols[i] = fmt.Sprintf("%s.%s", alias, f.Name)
	}
	return strings.Join(cols, ", ")
}

func (h *helper) ColumnBuilder(maker string, f *field) string {
	s := fmt.Sprintf(`%s.Col("%s", "%s")`, maker, f.Name, f.Column)
	if f.Tag.IsPK {
		s += ".PK()"
	}
	return s + ".Bld()"
}

type field struct {
	Name   string
	Column string
	Tag    ColumnTag
}

func GenerateCustomBuilders(opts Opts) error {
	conf := loader.Config{}
	if opts.IsTestFile {
		conf.ImportWithTests(opts.PkgPath)
	} else {
		conf.Import(opts.PkgPath)
	}
	prg, err := conf.Load()
	if err != nil {
		return nil
	}

	var tables types.Object
	var tablesPkg *types.Package

	pkg := prg.Package(opts.PkgPath)
	if opts.IsTestFile {
		testPkg := prg.Package(opts.PkgPath + "_test")
		tables, tablesPkg = findTablesDef(opts.TablesStructName, testPkg, pkg)
	} else {
		tables, tablesPkg = findTablesDef(opts.TablesStructName, pkg)
	}

	if tables == nil {
		return fmt.Errorf("Table definition struct '%s' not found", opts.TablesStructName)
	}

	tablesT, ok := tables.Type().Underlying().(*types.Struct)
	if !ok {
		return errors.Wrapf(err, "%s is not struct", opts.TablesStructName)
	}

	helpers := make([]*helper, tablesT.NumFields())

	for i := 0; i < tablesT.NumFields(); i++ {
		fld := tablesT.Field(i)

		tableName := fld.Name()

		var fldT *types.Struct
		fldVar, ok := fld.Type().(*types.Named)
		if ok {
			fldT, ok = fldVar.Underlying().(*types.Struct)
		}
		if !ok {
			return fmt.Errorf(
				"%s contains non struct field %s",
				opts.TablesStructName,
				tableName,
			)
		}

		modelPkgName := ""
		modelPkg := fldVar.Obj().Pkg()
		if modelPkg.Name() != tablesPkg.Name() {
			modelPkgName = modelPkg.Name()
		}

		tableTag, err := ParseTableTag(getTag(tablesT.Tag(i), "goq"))
		if err != nil {
			return errors.Wrapf(err, "failed to parse tag of Tables.%s", tableName)
		}

		helperName := tableTag.HelperName
		if helperName == "" {
			helperName = util.ColToFld(tableName)
		}

		modelName := fldVar.Obj().Name()
		fields, err := listColumnFields(modelName, fldT)
		if err != nil {
			return err
		}
		helpers[i] = &helper{
			Name:         helperName,
			TableName:    tableName,
			ModelPkgName: modelPkgName,
			ModelName:    modelName,
			Fields:       fields,
		}

	}

	src, err := execTemplate(tablesPkg.Name(), helpers, nil)
	if src != nil {
		file, err := createOutFile(opts.OutFile)
		if err != nil {
			return err
		}
		defer file.Close()
		file.Write(src)
	}

	return err
}

func findTablesDef(structName string, pkgs ...*loader.PackageInfo) (types.Object, *types.Package) {
	for _, pkg := range pkgs {
		if pkg != nil {
			if tables := pkg.Pkg.Scope().Lookup(structName); tables != nil {
				return tables, pkg.Pkg
			}
		}
	}
	return nil, nil
}

func listColumnFields(modelName string, modelT *types.Struct) ([]*field, error) {
	var fields []*field
	for i := 0; i < modelT.NumFields(); i++ {
		fld := modelT.Field(i)
		if !fld.Exported() {
			continue
		}

		tag, err := ParseColumnTag(getTag(modelT.Tag(i), "goq"))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse tag of %s.%s", modelName, fld.Name())
		}

		if tag.NotCol {
			continue
		}

		if fld.Anonymous() { // Embedded struct
			if ebdT, ok := fld.Type().Underlying().(*types.Struct); ok {
				ebdFields, err := listColumnFields(modelName, ebdT)
				if err != nil {
					return nil, err
				}
				fields = append(fields, ebdFields...)
			}
		} else {
			colName := tag.ColName
			if colName == "" {
				colName = util.FldToCol(fld.Name())
			}
			fields = append(fields, &field{
				Name:   fld.Name(),
				Column: colName,
				Tag:    tag,
			})
		}
	}

	return fields, nil
}

func createOutFile(outFile string) (*os.File, error) {
	_, err := os.Stat(outFile)

	if err == nil {
		err = os.Remove(outFile)
		if err != nil {
			return nil, fmt.Errorf("failed to remove %s", outFile)
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Create(outFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s", outFile)
	}
	return file, nil
}

func getTag(tag, key string) string {
	return reflect.StructTag(tag).Get(key)
}
