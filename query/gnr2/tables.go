package gnr2

// XXX: 面倒だけど、テーブルは個別に定義するしかなさそう。

type Joinner struct {
	table Table
	on    PredExpr
}

func (jd *Joinner) Inner() JoinOn {
	return JoinOn{jd.table, jd.on, JOIN_INNER}
}

type UsersTable struct {
	*SliceCollectorMaker
	empModel interface{}
	name     string
	alias    string

	ID   Column
	Name Column
}

func (t *UsersTable) Posts(t2 *PostsTable) *Joinner {
	return &Joinner{t2, t.ID.Eq(t2.UserID)}
}

func (t *UsersTable) TableName() string  { return t.name }
func (t *UsersTable) TableAlias() string { return t.alias }

func (t *UsersTable) All() ExprListExpr {
	cols := t.Columns()
	exps := make([]Expr, len(cols))
	for i, c := range cols {
		exps[i] = c
	}
	return &exprListExpr{exps}
}

func (t *UsersTable) Columns() []Column {
	return []Column{t.ID, t.Name}
}

// ちょっとコスト高すぎ..?
func (t *UsersTable) As(alias string) *UsersTable {
	origCols := t.Columns()
	cols := make([]Column, len(origCols))
	for i, c := range origCols {
		col := column{
			tableAlias: alias,
			tableName:  c.TableName(),
			structName: c.StructName(),
			name:       c.ColumnName(),
			fieldName:  c.FieldName(),
		}.init()
		cols[i] = &col
	}

	t2 := *t
	t2.alias = alias
	t2.ID = cols[0]
	t2.Name = cols[1]
	t2.SliceCollectorMaker = NewSliceCollectorMaker(t.empModel, cols, alias)
	return &t2
}

type PostsTable struct {
	name  string
	alias string

	ID     Column
	UserID Column
}

func (t *PostsTable) TableName() string  { return t.name }
func (t *PostsTable) TableAlias() string { return t.alias }

func (t *PostsTable) All() ExprListExpr {
	cols := t.Columns()
	exps := make([]Expr, len(cols))
	for i, c := range cols {
		exps[i] = c
	}
	return &exprListExpr{exps}
}

func (t *PostsTable) Columns() []Column {
	return []Column{t.ID, t.UserID}
}

func (t *PostsTable) As(alias string) *PostsTable {
	origCols := t.Columns()
	cols := make([]Column, len(origCols))
	for i, c := range origCols {
		col := column{
			tableAlias: alias,
			tableName:  c.TableName(),
			structName: c.StructName(),
			name:       c.ColumnName(),
			fieldName:  c.FieldName(),
		}.init()
		cols[i] = &col
	}

	t2 := *t
	t2.alias = alias
	t2.ID = cols[0]
	t2.UserID = cols[1]
	// t2.SliceCollectorMaker = NewSliceCollectorMaker(User{}, cols, alias)
	return &t2
}
