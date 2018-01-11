package gnr2

// こういう感じのコードを generate する。

func all(cols []Column) ExprListExpr {
	exps := make([]Expr, len(cols))
	for i, c := range cols {
		exps[i] = c
	}
	return &exprListExpr{exps}
}

type UsersTable struct {
	*CollectorMaker
	model interface{}
	name  string
	alias string

	ID   Column
	Name Column
}

func (t *UsersTable) TableName() string  { return t.name }
func (t *UsersTable) TableAlias() string { return t.alias }
func (t *UsersTable) All() ExprListExpr  { return all(t.Columns()) }

func (t *UsersTable) Columns() []Column {
	return []Column{t.ID, t.Name}
}

func (t *UsersTable) Posts(t2 *PostsTable) *JoinDef {
	return &JoinDef{t2, t.ID.Eq(t2.UserID)}
}

// ちょっとコスト高すぎ..?
func (t *UsersTable) As(alias string) *UsersTable {
	t2 := *t
	t2.alias = alias
	copyTableAs(alias, t, &t2, t.model)
	return &t2
}

type PostsTable struct {
	*CollectorMaker
	model interface{}
	name  string
	alias string

	ID     Column
	UserID Column
}

func (t *PostsTable) TableName() string  { return t.name }
func (t *PostsTable) TableAlias() string { return t.alias }
func (t *PostsTable) All() ExprListExpr  { return all(t.Columns()) }

func (t *PostsTable) Columns() []Column {
	return []Column{t.ID, t.UserID}
}

func (t *PostsTable) As(alias string) *PostsTable {
	t2 := *t
	t2.alias = alias
	copyTableAs(alias, t, &t2, t.model)
	return &t2
}

type PrefsTable struct {
	*CollectorMaker
	model interface{}
	name  string
	alias string

	ID   Column
	Name Column
}

func (t *PrefsTable) TableName() string  { return t.name }
func (t *PrefsTable) TableAlias() string { return t.alias }
func (t *PrefsTable) All() ExprListExpr  { return all(t.Columns()) }

func (t *PrefsTable) Columns() []Column {
	return []Column{t.ID, t.Name}
}

func (t *PrefsTable) As(alias string) *PrefsTable {
	t2 := *t
	t2.alias = alias
	copyTableAs(alias, t, &t2, t.model)
	return &t2
}

type CitiesTable struct {
	*CollectorMaker
	model interface{}
	name  string
	alias string

	ID     Column
	Name   Column
	PrefID Column
}

func (t *CitiesTable) TableName() string  { return t.name }
func (t *CitiesTable) TableAlias() string { return t.alias }
func (t *CitiesTable) All() ExprListExpr  { return all(t.Columns()) }

func (t *CitiesTable) Columns() []Column {
	return []Column{t.ID, t.Name, t.PrefID}
}

func (t *CitiesTable) As(alias string) *CitiesTable {
	t2 := *t
	t2.alias = alias
	copyTableAs(alias, t, &t2, t.model)
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
		model: User{},
		name:  "users",
		ID:    cm.Col("ID", "id"),
		Name:  cm.Col("Name", "name"),
	}
	g.Users.CollectorMaker = NewCollectorMaker(g.Users.model, g.Users.Columns(), "")

	cm = ColumnMaker{"posts", "", "Post"}
	g.Posts = &PostsTable{
		model:  Post{},
		name:   "posts",
		ID:     cm.Col("ID", "id"),
		UserID: cm.Col("UserID", "user_id"),
	}
	// g.Posts.CollectorMaker = NewCollectorMaker(g.Posts.model, g.Posts.Columns(), "")

	cm = ColumnMaker{"prefectures", "", "Pref"}
	g.Prefs = &PrefsTable{
		model: Pref{},
		name:  "prefectures",
		ID:    cm.Col("ID", "id"),
		Name:  cm.Col("Name", "name"),
	}
	g.Prefs.CollectorMaker = NewCollectorMaker(g.Prefs.model, g.Prefs.Columns(), "")

	cm = ColumnMaker{"cities", "", "City"}
	g.Cities = &CitiesTable{
		model:  City{},
		name:   "cities",
		ID:     cm.Col("ID", "id"),
		Name:   cm.Col("Name", "name"),
		PrefID: cm.Col("PrefID", "prefecture_id"),
	}
	g.Cities.CollectorMaker = NewCollectorMaker(g.Cities.model, g.Cities.Columns(), "")

	return &g
}
