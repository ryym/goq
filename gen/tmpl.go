package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"text/template"

	"github.com/pkg/errors"
)

func writeTemplate(w io.Writer, pkgName string, helpers []*helper) error {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "package %s\n\nimport \"github.com/ryym/goq/gql\"\n", pkgName)

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

const tableTmpl = `
type {{.Name}} struct {
	model {{.ModelName}}
	name  string
	alias string
	{{range .Fields}}
	{{.Name}} gql.Column{{end}}
}

func New{{.Name}}() *{{.Name}} {
	cm := gql.NewColumnMaker("{{.TableName}}", "{{.ModelName}}")
	return &{{.Name}}{
		model: {{.ModelName}}{},
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
