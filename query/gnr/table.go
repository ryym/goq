package gnr

import q "github.com/ryym/goq/query"

// XXX: 面倒だけど、テーブルは個別に定義するしかなさそう。
// columnExpr などは後々 Column などに変える。

type Joinner struct {
	on q.PredExpr
}

func (jd *Joinner) Inner() q.JoinOn {
	return q.JoinOn{jd.on, q.JOIN_INNER}
}

type UsersTable struct {
	name string

	ID   q.Expr
	Name q.Expr
}

func (t *UsersTable) Posts(t2 *PostsTable) Joinner {
	return Joinner{t.ID.Eq(t2.UserID)}
}

func (t *UsersTable) TableName() string { return t.name }
func (t *UsersTable) Alias() string     { return "" }

func (u *UsersTable) All() q.ExprListExpr {
	return &exprListExpr{[]q.Queryable{u.ID, u.Name}}
}

func (u *UsersTable) As(alias string) *UsersAliased {
	origCols := u.All().(*exprListExpr).qs
	cols := make([]columnExpr, len(origCols))
	for i, c := range origCols {
		// cols[i] = *c.(*columnExpr)
		cols[i] = *c.(*Ops).Queryable.(*columnExpr)
		cols[i].tableAlias = alias
	}
	u2 := *u
	u2.ID = &Ops{&cols[0]}
	u2.Name = &Ops{&cols[1]}
	return &UsersAliased{u2, alias}
}

type UsersAliased struct {
	UsersTable
	alias string
}

func (t *UsersAliased) Alias() string { return t.alias }

type PostsTable struct {
	name string

	ID     q.Expr
	UserID q.Expr
}

func (t *PostsTable) TableName() string { return t.name }
func (t *PostsTable) Alias() string     { return "" }

func (u *PostsTable) All() q.ExprListExpr {
	return &exprListExpr{[]q.Queryable{u.ID, u.UserID}}
}

func (u *PostsTable) As(alias string) *PostsAliased {
	origCols := u.All().(*exprListExpr).qs
	cols := make([]columnExpr, len(origCols))
	for i, c := range origCols {
		// cols[i] = *c.(*columnExpr)
		cols[i] = *c.(*Ops).Queryable.(*columnExpr)
		cols[i].tableAlias = alias
	}
	u2 := *u
	u2.ID = &Ops{&cols[0]}
	u2.UserID = &Ops{&cols[1]}
	return &PostsAliased{u2, alias}
}

type PostsAliased struct {
	PostsTable
	alias string
}

func (t *PostsAliased) Alias() string { return t.alias }
