package gen

import (
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"
)

func writeTemplate(w io.Writer, pkgName string, helpers []*helper) error {
	tableT := template.Must(template.New("table").Parse(tableTmpl))

	fmt.Fprintf(w, "package %s\n\nimport \"github.com/ryym/goq/gql\"\n", pkgName)

	var err error
	for _, h := range helpers {
		err = tableT.Execute(w, h)
		if err != nil {
			return errors.Wrapf(err, "failed to create template of %s", h.Name)
		}
	}

	// TODO: Need to apply `gofmt` to the generated Go file.

	return nil
}

const tableTmpl = `
type {{.Name}} struct {
	model {{.ModelName}}
	name  string
	alias string
	{{range .Fields}}
	{{.Name}} gql.Column {{end}}
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
