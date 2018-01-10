package gnr2

// こういう感じのコードを generate する。

type UsersTable struct {
	*CollectorMaker
	empModel interface{}
	name     string
	alias    string

	ID   Column
	Name Column
}

func (t *UsersTable) TableName() string  { return t.name }
func (t *UsersTable) TableAlias() string { return t.alias }

func (t *UsersTable) Columns() []Column {
	return []Column{t.ID, t.Name}
}

func (t *UsersTable) All() ExprListExpr {
	cols := t.Columns()
	exps := make([]Expr, len(cols))
	for i, c := range cols {
		exps[i] = c
	}
	return &exprListExpr{exps}
}

func (t *UsersTable) Posts(t2 *PostsTable) *JoinDef {
	return &JoinDef{t2, t.ID.Eq(t2.UserID)}
}

// ちょっとコスト高すぎ..?
func (t *UsersTable) As(alias string) *UsersTable {
	t2 := *t
	t2.alias = alias
	t2.CollectorMaker = NewCollectorMaker(t.empModel, t2.Columns(), alias)
	copyTableAs(alias, t, &t2)
	return &t2
}

type PostsTable struct {
	*CollectorMaker
	empModel interface{}
	name     string
	alias    string

	ID     Column
	UserID Column
}

func (t *PostsTable) TableName() string  { return t.name }
func (t *PostsTable) TableAlias() string { return t.alias }

func (t *PostsTable) Columns() []Column {
	return []Column{t.ID, t.UserID}
}

func (t *PostsTable) All() ExprListExpr {
	cols := t.Columns()
	exps := make([]Expr, len(cols))
	for i, c := range cols {
		exps[i] = c
	}
	return &exprListExpr{exps}
}

func (t *PostsTable) As(alias string) *PostsTable {
	t2 := *t
	t2.alias = alias
	t2.CollectorMaker = NewCollectorMaker(t.empModel, t2.Columns(), alias)
	copyTableAs(alias, t, &t2)
	return &t2
}

type PrefsTable struct {
	*CollectorMaker
	empModel interface{}
	name     string
	alias    string

	ID   Column
	Name Column
}

func (t *PrefsTable) TableName() string  { return t.name }
func (t *PrefsTable) TableAlias() string { return t.alias }

func (t *PrefsTable) Columns() []Column {
	return []Column{t.ID, t.Name}
}

func (t *PrefsTable) All() ExprListExpr {
	cols := t.Columns()
	exps := make([]Expr, len(cols))
	for i, c := range cols {
		exps[i] = c
	}
	return &exprListExpr{exps}
}

func (t *PrefsTable) As(alias string) *PrefsTable {
	t2 := *t
	t2.alias = alias
	t2.CollectorMaker = NewCollectorMaker(t.empModel, t2.Columns(), alias)
	copyTableAs(alias, t, &t2)
	return &t2
}

type CitiesTable struct {
	*CollectorMaker
	empModel interface{}
	name     string
	alias    string

	ID     Column
	Name   Column
	PrefID Column
}

func (t *CitiesTable) TableName() string  { return t.name }
func (t *CitiesTable) TableAlias() string { return t.alias }

func (t *CitiesTable) Columns() []Column {
	return []Column{t.ID, t.Name, t.PrefID}
}

func (t *CitiesTable) All() ExprListExpr {
	cols := t.Columns()
	exps := make([]Expr, len(cols))
	for i, c := range cols {
		exps[i] = c
	}
	return &exprListExpr{exps}
}

func (t *CitiesTable) As(alias string) *CitiesTable {
	t2 := *t
	t2.alias = alias
	t2.CollectorMaker = NewCollectorMaker(t.empModel, t2.Columns(), alias)
	copyTableAs(alias, t, &t2)
	return &t2
}

type GQL struct {
	*Goq   // 実際には QueryBuilder interface
	Users  *UsersTable
	Posts  *PostsTable
	Prefs  *PrefsTable
	Cities *CitiesTable
}

func NewGQL() *GQL {
	g := GQL{Goq: &Goq{}}

	cm := ColumnMaker{"users", "", "User"}
	g.Users = &UsersTable{
		empModel: User{},
		name:     "users",
		ID:       cm.Col("ID", "id"),
		Name:     cm.Col("Name", "name"),
	}
	g.Users.CollectorMaker = NewCollectorMaker(g.Users.empModel, g.Users.Columns(), "")

	cm = ColumnMaker{"posts", "", "Post"}
	g.Posts = &PostsTable{
		empModel: Post{},
		name:     "posts",
		ID:       cm.Col("ID", "id"),
		UserID:   cm.Col("UserID", "user_id"),
	}
	// g.Posts.CollectorMaker = NewCollectorMaker(g.Posts.empModel, g.Posts.Columns(), "")

	cm = ColumnMaker{"prefectures", "", "Pref"}
	g.Prefs = &PrefsTable{
		empModel: Pref{},
		name:     "prefectures",
		ID:       cm.Col("ID", "id"),
		Name:     cm.Col("Name", "name"),
	}
	g.Prefs.CollectorMaker = NewCollectorMaker(g.Prefs.empModel, g.Prefs.Columns(), "")

	cm = ColumnMaker{"cities", "", "City"}
	g.Cities = &CitiesTable{
		empModel: City{},
		name:     "cities",
		ID:       cm.Col("ID", "id"),
		Name:     cm.Col("Name", "name"),
		PrefID:   cm.Col("PrefID", "prefecture_id"),
	}
	g.Cities.CollectorMaker = NewCollectorMaker(g.Cities.empModel, g.Cities.Columns(), "")

	return &g
}
