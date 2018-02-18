package gql

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
)

type testDialect struct{}

func (d *testDialect) Placeholder(typ string, prevArgs []interface{}) string {
	ph := fmt.Sprintf("$%d", len(prevArgs)+1)
	if typ != "" {
		ph += "::" + typ
	}
	return ph
}

func (d *testDialect) QuoteIdent(v string) string {
	return fmt.Sprintf("`%s`", v)
}

type usersTable struct {
	Table
	name  string
	alias string
	ID    *Column
	Name  *Column
}

func (t *usersTable) As(alias string) *usersTable { return t /* FOR NOW */ }

func TestBasicExprs(t *testing.T) {
	z := NewBuilder(&testDialect{})
	cm := NewColumnMaker("User", "users")
	ID := cm.Col("ID", "id").Bld()
	Name := cm.Col("Name", "name").Bld()
	Users := &usersTable{
		Table: Table{"users", "", []*Column{ID, Name}},
		ID:    ID,
		Name:  Name,
	}

	type user struct {
		ID   int
		Name string
	}

	var tests = []struct {
		gql  QueryApplier
		sql  string
		args []interface{}
	}{
		{
			gql:  z.Var(1),
			sql:  "$1",
			args: []interface{}{1},
		},
		{
			gql:  z.Var(1).As("one"),
			sql:  "$1 AS `one`",
			args: []interface{}{1},
		},
		{
			gql:  z.Var(1).Eq(2),
			sql:  "$1 = $2",
			args: []interface{}{1, 2},
		},
		{
			gql:  z.Var(1).Eq(2).As("t"),
			sql:  "$1 = $2 AS `t`",
			args: []interface{}{1, 2},
		},
		{
			gql:  ID.Between(0, 5),
			sql:  "`users`.`id` BETWEEN $1 AND $2",
			args: []interface{}{0, 5},
		},
		{
			gql:  z.Parens(ID.Add(4).Sbt(3)).Mlt(2).Dvd(1),
			sql:  "(`users`.`id` + $1 - $2) * $3 / $4",
			args: []interface{}{4, 3, 2, 1},
		},
		{
			gql:  Name.Concat("a").Concat(Name),
			sql:  "`users`.`name` || $1 || `users`.`name`",
			args: []interface{}{"a"},
		},
		{
			gql: z.And(
				ID.Eq(1),
				z.Or(ID.Lte(3), z.Var(1).Gt(ID)),
				Name.IsNotNull(),
			),
			sql:  "(`users`.`id` = $1 AND (`users`.`id` <= $2 OR $3 > `users`.`id`) AND `users`.`name` IS NOT NULL)",
			args: []interface{}{1, 3, 1},
		},
		{
			gql:  z.Func("foo", 1, ID, 2).Add(3),
			sql:  "foo($1, `users`.`id`, $2) + $3",
			args: []interface{}{1, 2, 3},
		},
		{
			gql:  z.Coalesce(Name, z.Var(20)),
			sql:  "COALESCE(`users`.`name`, $1)",
			args: []interface{}{20},
		},
		{
			gql: z.Case(
				z.When(ID.Eq(1)).Then(-1),
				z.When(ID.Eq(2)).Then(-3),
			).Else(0).Add(1),
			sql: "CASE WHEN `users`.`id` = $1 THEN $2 WHEN `users`.`id` = $3 THEN $4" +
				" ELSE $5 END + $6",
			args: []interface{}{1, -1, 2, -3, 0, 1},
		},
		{
			gql:  z.Select(ID, Name, z.Var(1).Add(ID).As("test")),
			sql:  "SELECT `users`.`id`, `users`.`name`, $1 + `users`.`id` AS `test`",
			args: []interface{}{1},
		},
		{
			gql: z.Select(Users.All()).From(Users).Joins(
				z.LeftJoin(Users).On(Name.Eq("bob")),
			).Where(
				ID.Gte(3),
				Name.Like("%bob"),
			).GroupBy(
				ID,
				Name,
			).Having(
				z.Count(ID).Lt(100),
			).OrderBy(ID).Limit(10).Offset(20),
			sql: "SELECT `users`.`id`, `users`.`name` FROM `users` " +
				"LEFT OUTER JOIN `users` ON `users`.`name` = $1 " +
				"WHERE (`users`.`id` >= $2 AND `users`.`name` LIKE $3) " +
				"GROUP BY `users`.`id`, `users`.`name` " +
				"HAVING (COUNT(`users`.`id`) < $4) " +
				"ORDER BY `users`.`id` LIMIT 10 OFFSET 20",
			args: []interface{}{"bob", 3, "%bob", 100},
		},
		{
			gql:  z.Select(z.Select(z.Var(1)).As("subquery")),
			sql:  "SELECT (SELECT $1) AS `subquery`",
			args: []interface{}{1},
		},
		{
			gql: z.Select(z.Var(1)).From(
				z.Select(z.Var(3)).As("subquery"),
			).Joins(
				z.RightJoin(z.Select(Users.ID).From(Users).As("u")).On(
					z.Var(5).Eq(7),
				),
			),
			sql: "SELECT $1 FROM (SELECT $2) AS `subquery` " +
				"RIGHT OUTER JOIN (SELECT `users`.`id` FROM `users`) AS `u` ON $3 = $4",
			args: []interface{}{1, 3, 5, 7},
		},
		{
			gql: z.Select(
				z.Col("", "id"), z.Col("f", "title"), z.Col("foo", "body").As("content"),
				z.Col("", "count").Add(3),
			).From(
				z.Table("foo").As("f"),
			),
			sql:  "SELECT `id`, `f`.`title`, `foo`.`body` AS `content`, `count` + $1 FROM `foo` AS `f`",
			args: []interface{}{3},
		},
		{
			gql:  z.Select(z.Var(1)).From(Users).Joins(&Join{Users, Users.ID.Eq(3)}),
			sql:  "SELECT $1 FROM `users` INNER JOIN `users` ON `users`.`id` = $2",
			args: []interface{}{1, 3},
		},
		{
			gql:  z.Select(z.Var(1).As("num")).OrderBy(z.Name("num")),
			sql:  "SELECT $1 AS `num` ORDER BY `num`",
			args: []interface{}{1},
		},
		{
			gql:  z.Select(z.Var(1)).OrderBy(Name.Asc(), ID.Desc(), ID),
			sql:  "SELECT $1 ORDER BY `users`.`name`, `users`.`id` DESC, `users`.`id`",
			args: []interface{}{1},
		},
		{
			// Postgres needs type information for placeholders in some cases.
			gql:  z.Select(z.VarT(1, "int")),
			sql:  "SELECT $1::int",
			args: []interface{}{1},
		},
		{
			gql:  z.Select(z.Null().As("a")).From(Users).Where(z.Null().Eq(z.Null())),
			sql:  "SELECT NULL AS `a` FROM `users` WHERE (NULL = NULL)",
			args: nil,
		},
		{
			gql:  z.Select(Users.Except(Users.ID), z.Null()).From(Users),
			sql:  "SELECT `users`.`name`, NULL FROM `users`",
			args: nil,
		},
		{
			gql: z.InsertInto(Users).ValuesMap(Values{
				Users.ID:   1,
				Users.Name: "bob",
			}),
			sql:  "INSERT INTO `users` VALUES ($1, $2)",
			args: []interface{}{1, "bob"},
		},
		{
			gql:  z.InsertInto(Users).Values(user{1, "bob"}),
			sql:  "INSERT INTO `users` VALUES ($1, $2)",
			args: []interface{}{1, "bob"},
		},
		{
			gql:  z.InsertInto(Users, Users.Name).Values(user{1, "bob"}),
			sql:  "INSERT INTO `users` (`name`) VALUES ($1)",
			args: []interface{}{"bob"},
		},
		{
			gql: z.InsertInto(
				Users,
				Users.Except(Users.Name).Columns()...,
			).Values(user{1, "bob"}),
			sql:  "INSERT INTO `users` (`id`) VALUES ($1)",
			args: []interface{}{1},
		},
	}

	for i, test := range tests {
		q := z.Query(test.gql)
		if query := q.Query(); query != test.sql {
			t.Errorf("[%d] Query diff\nGOT : %s\nWANT: %s", i, query, test.sql)
		}
		if diff := deep.Equal(q.Args(), test.args); diff != nil {
			t.Errorf("[%d] Args diff\n%s", i, diff)
		}
	}
}
