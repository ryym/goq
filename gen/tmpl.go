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
	modelPkgs map[string]bool,
) error {
	buf := new(bytes.Buffer)
	buf.Write([]byte(fmt.Sprintf("package %s\n\n", pkgName)))

	writeImports(buf, modelPkgs)
	tableT := template.Must(template.New("table").Parse(tableTmpl))

	var err error
	for _, h := range helpers {
		err = tableT.Execute(buf, h)
		if err != nil {
			return errors.Wrapf(err, "failed to create template of %s", h.Name)
		}
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	w.Write(src)

	return nil
}

func writeImports(buf io.Writer, modelPkgs map[string]bool) {
	if len(modelPkgs) == 0 {
		buf.Write([]byte("\n\nimport \"github.com/ryym/goq/gql\"\n"))
	} else {
		buf.Write([]byte("import (\n"))

		paths := make([]string, len(modelPkgs)+1)
		paths[0] = "github.com/ryym/goq/gql"
		i := 1
		for path, _ := range modelPkgs {
			paths[i] = path
			i++
		}
		sort.Strings(paths)
		for _, path := range paths {
			fmt.Fprintf(buf, "\"%s\"\n", path)
		}

		buf.Write([]byte(")\n"))
	}

}

const tableTmpl = `
type {{.Name}} struct {
	model {{.ModelFullName}}
	name  string
	alias string
	{{range .Fields}}
	{{.Name}} gql.Column{{end}}
}

func New{{.Name}}() *{{.Name}} {
	cm := gql.NewColumnMaker("{{.TableName}}", "{{.ModelName}}")
	return &{{.Name}}{
		model: {{.ModelFullName}}{},
		name:  "{{.TableName}}",
		{{range .Fields}}
		{{.Name}}: cm.Col("{{.Name}}", "{{.Column}}"),{{end}}
	}
}

func (t *{{.Name}}) TableName() string     { return t.name }
func (t *{{.Name}}) TableAlias() string    { return t.alias }
func (t *{{.Name}}) All() gql.ExprListExpr { return gql.AllCols(t.Columns()) }
func (t *{{.Name}}) Columns() []gql.Column {
	return []gql.Column{ {{.JoinFields "t"}} }
}
func (t *{{.Name}}) As(alias string) *{{.Name}} {
	t2 := *t
	t2.alias = alias
	gql.CopyTableAs(alias, t, &t2)
	return &t2
}
`
