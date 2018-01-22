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
	Pkg              string
	OutFile          string
	TablesStructName string
}

type helper struct {
	Name         string
	TableName    string
	ModelPkgName string
	ModelName    string
	Fields       []*field
}

func (h *helper) ModelFullName() string {
	if h.ModelPkgName != "" {
		return fmt.Sprintf("%s.%s", h.ModelPkgName, h.ModelName)
	}
	return h.ModelName
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
	conf := loader.Config{}
	conf.Import(opts.Pkg)
	prg, err := conf.Load()
	if err != nil {
		return nil
	}

	pkg := prg.Package(opts.Pkg)
	if pkg == nil {
		return fmt.Errorf("package %s not found", opts.Pkg)
	}

	tablesPkg := pkg.Pkg
	tables := tablesPkg.Scope().Lookup(opts.TablesStructName)
	if tables == nil {
		return fmt.Errorf("struct %s not found", opts.TablesStructName)
	}

	tablesT, ok := tables.Type().Underlying().(*types.Struct)
	if !ok {
		return errors.Wrapf(err, "%s is not struct", opts.TablesStructName)
	}

	helpers := make([]*helper, tablesT.NumFields())

	for i := 0; i < tablesT.NumFields(); i++ {
		fld := tablesT.Field(i)
		tableName := fld.Name()
		fldVar := fld.Type().(*types.Named)
		fldT, ok := fldVar.Underlying().(*types.Struct)
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

	file, err := createOutFile(opts.OutFile)
	if err != nil {
		return err
	}
	defer file.Close() // TODO: Remove file on error

	err = writeTemplate(file, tablesPkg.Name(), helpers, nil)
	if err != nil {
		return err
	}

	return nil
}

func listColumnFields(modelName string, modelT *types.Struct) ([]*field, error) {
	var fields []*field
	for i := 0; i < modelT.NumFields(); i++ {
		fld := modelT.Field(i)
		if !fld.Exported() {
			continue
		}

		tag, err := ParseModelTag(getTag(modelT.Tag(i), "goq"))
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
			})
		}
	}

	return fields, nil
}

func createOutFile(outPath string) (*os.File, error) {
	if _, err := os.Stat(outPath); err == nil {
		err = os.Remove(outPath)
		if err != nil {
			return nil, fmt.Errorf("failed to remove %s", outPath)
		}
	}

	file, err := os.Create(outPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s", outPath)
	}
	return file, nil
}

func getTag(tag, key string) string {
	return reflect.StructTag(tag).Get(key)
}
