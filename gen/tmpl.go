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

func writeTemplate(
	w io.Writer,
	pkgName string,
	helpers []*helper,
	pkgs map[string]bool,
) error {
	buf := new(bytes.Buffer)
	buf.Write([]byte(fmt.Sprintf("package %s\n\n", pkgName)))

	writeImports(buf, pkgs)
	tableT := template.Must(template.New("table").Parse(tableTmpl))

	var err error
	for _, h := range helpers {
		err = tableT.Execute(buf, h)
		if err != nil {
			return errors.Wrapf(err, "failed to execute template of %s", h.Name)
		}
	}

	gqlStructT := template.Must(template.New("gqlStruct").Parse(gqlStructTmpl))
	err = gqlStructT.Execute(buf, helpers)
	if err != nil {
		return errors.Wrap(err, "failed to execute Builder struct template")
	}

	newBuilderT := template.Must(template.New("newBuilder").Parse(newBuilderTmpl))
	err = newBuilderT.Execute(buf, helpers)
	if err != nil {
		return errors.Wrap(err, "failed to execute NewBuilder() template")
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	w.Write(src)

	return nil
}

func writeImports(buf io.Writer, pkgs map[string]bool) {
	buf.Write([]byte("import (\n"))

	paths := []string{
		"github.com/ryym/goq",
		"github.com/ryym/goq/cllct",
		"github.com/ryym/goq/dialect",
		"github.com/ryym/goq/gql",
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
	gql.Table
	*cllct.ModelCollectorMaker
	{{range .Fields}}
	{{.Name}} *gql.Column{{end}}
}

func New{{.Name}}(alias string) *{{.Name}} {
	cm := gql.NewColumnMaker("{{.ModelName}}", "{{.TableName}}").As(alias)
	t := &{{.Name}}{
		Table: gql.NewTable("{{.TableName}}", alias),
		{{range .Fields}}
		{{.Name}}: cm.Col("{{.Name}}", "{{.Column}}"),{{end}}
	}
	t.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), alias)
	return t
}

func (t *{{.Name}}) As(alias string) *{{.Name}} { return New{{.Name}}(alias) }
func (t *{{.Name}}) All() gql.ExprListExpr      { return gql.AllCols(t.Columns()) }
func (t *{{.Name}}) Columns() []*gql.Column {
	return []*gql.Column{ {{.JoinFields "t"}} }
}
`

const gqlStructTmpl = `
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
