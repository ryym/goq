package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"sort"
	"text/template"

	"github.com/pkg/errors"
)

func execTemplate(
	pkgName string,
	helpers []*helper,
	pkgs map[string]bool,
) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Write([]byte(fmt.Sprintf(`
// DO NOT EDIT. This code is generated by Goq.
// https://github.com/ryym/goq

package %s
`, pkgName)))

	writeImports(buf, helpers, pkgs)
	tableT := template.Must(template.New("table").Parse(tableTmpl))

	var err error
	for _, h := range helpers {
		err = tableT.Execute(buf, h)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to execute template of %s", h.Name)
		}
	}

	builderT := template.Must(template.New("builder").Parse(builderTmpl))
	err = builderT.Execute(buf, helpers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute Builder struct template")
	}

	newBuilderT := template.Must(template.New("newBuilder").Parse(newBuilderTmpl))
	err = newBuilderT.Execute(buf, helpers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute NewBuilder() template")
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		err = errors.Wrap(err, "failed to format generated code")
	}
	return src, err
}

func writeImports(buf io.Writer, helpers []*helper, pkgs map[string]bool) {
	buf.Write([]byte("import (\n"))

	paths := []string{
		"github.com/ryym/goq",
		"github.com/ryym/goq/dialect",
	}
	if len(helpers) > 0 {
		paths = append(paths, "github.com/ryym/goq/goql")
	}

	for path, _ := range pkgs {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Fprintf(buf, "\"%s\"\n", path)
	}

	buf.Write([]byte(")\n"))
}

const tableTmpl = `
type {{.Name}} struct {
	goql.Table
	*goq.ModelCollectorMaker
	{{range .Fields}}
	{{.Name}} *goql.Column{{end}}
}

func New{{.Name}}(alias string) *{{.Name}} {
	{{if .Fields}}cm := goql.NewColumnMaker("{{.ModelName}}", "{{.TableName}}").As(alias){{end}}
	t := &{{.Name}}{
		{{range .Fields}}
		{{.Name}}: {{$.ColumnBuilder "cm" .}},{{end}}
	}
	cols := []*goql.Column{ {{.JoinFields "t"}} }
	t.Table = goql.NewTable("{{.TableName}}", alias, cols)
	t.ModelCollectorMaker = goq.NewModelCollectorMaker(cols, alias)
	return t
}

func (t *{{.Name}}) As(alias string) *{{.Name}} { return New{{.Name}}(alias) }
`

const builderTmpl = `
type Builder struct {
	*goq.Builder
	{{range .}}
	{{.Name}} *{{.Name}}{{end}}
}
`

const newBuilderTmpl = `
func NewBuilder(dl dialect.Dialect) *Builder {
	return &Builder{
		Builder: goq.NewBuilder(dl),
		{{range .}}
		{{.Name}}: New{{.Name}}(""),{{end}}
	}
}`
