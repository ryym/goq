package goql

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
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
	ID := cm.Col("ID", "id").PK().Bld()
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
		goql QueryApplier
		sql  string
		args []interface{}
	}{
		{
			goql: z.Var(1),
			sql:  "$1",
			args: []interface{}{1},
		},
		{
			goql: z.Var(1).As("one"),
			sql:  "$1 AS `one`",
			args: []interface{}{1},
		},
		{
			goql: z.Var(1).Eq(2),
			sql:  "$1 = $2",
			args: []interface{}{1, 2},
		},
		{
			goql: z.Var(1).Eq(2).As("t"),
			sql:  "$1 = $2 AS `t`",
			args: []interface{}{1, 2},
		},
		{
			goql: ID.Between(0, 5),
			sql:  "`users`.`id` BETWEEN $1 AND $2",
			args: []interface{}{0, 5},
		},
		{
			goql: ID.In([]int{1, 2, 3}),
			sql:  "`users`.`id` IN ($1, $2, $3)",
			args: []interface{}{1, 2, 3},
		},
		{
			goql: ID.NotIn([]int{1, 2, 3}),
			sql:  "`users`.`id` NOT IN ($1, $2, $3)",
			args: []interface{}{1, 2, 3},
		},
		{
			goql: ID.In(z.Select(z.Var(1))),
			sql:  "`users`.`id` IN (SELECT $1)",
			args: []interface{}{1},
		},
		{
			goql: z.Parens(ID.Add(4).Sbt(3)).Mlt(2).Dvd(1),
			sql:  "(`users`.`id` + $1 - $2) * $3 / $4",
			args: []interface{}{4, 3, 2, 1},
		},
		{
			goql: Name.Concat("a").Concat(Name),
			sql:  "`users`.`name` || $1 || `users`.`name`",
			args: []interface{}{"a"},
		},
		{
			// An alias is ignored.
			goql: z.Var(1).Add(z.Var(3).As("foo")),
			sql:  "$1 + $2",
			args: []interface{}{1, 3},
		},
		{
			goql: z.And(
				ID.Eq(1),
				z.Or(ID.Lte(3), z.Var(1).Gt(ID)),
				Name.IsNotNull(),
			),
			sql:  "(`users`.`id` = $1 AND (`users`.`id` <= $2 OR $3 > `users`.`id`) AND `users`.`name` IS NOT NULL)",
			args: []interface{}{1, 3, 1},
		},
		{
			goql: z.Func("foo", 1, ID, 2).Add(3),
			sql:  "foo($1, `users`.`id`, $2) + $3",
			args: []interface{}{1, 2, 3},
		},
		{
			goql: z.Coalesce(Name, z.Var(20)),
			sql:  "COALESCE(`users`.`name`, $1)",
			args: []interface{}{20},
		},
		{
			goql: z.Case(
				z.When(ID.Eq(1)).Then(-1),
				z.When(ID.Eq(2)).Then(-3),
			).Else(0).Add(1),
			sql: "CASE WHEN `users`.`id` = $1 THEN $2 WHEN `users`.`id` = $3 THEN $4" +
				" ELSE $5 END + $6",
			args: []interface{}{1, -1, 2, -3, 0, 1},
		},
		{
			goql: z.CaseOf(ID,
				z.When(1).Then(-1),
				z.When(2).Then(-3),
			).Else(0).Add(1),
			sql: "CASE `users`.`id` WHEN $1 THEN $2 WHEN $3 THEN $4" +
				" ELSE $5 END + $6",
			args: []interface{}{1, -1, 2, -3, 0, 1},
		},
		{
			goql: z.Select(ID, Name, z.Var(1).Add(ID).As("test")),
			sql:  "SELECT `users`.`id`, `users`.`name`, $1 + `users`.`id` AS `test`",
			args: []interface{}{1},
		},
		{
			goql: z.SelectDistinct(ID),
			sql:  "SELECT DISTINCT `users`.`id`",
			args: nil,
		},
		{
			goql: z.Select(Users.All()).From(Users).Joins(
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
				"WHERE `users`.`id` >= $2 AND `users`.`name` LIKE $3 " +
				"GROUP BY `users`.`id`, `users`.`name` " +
				"HAVING (COUNT(`users`.`id`) < $4) " +
				"ORDER BY `users`.`id` LIMIT 10 OFFSET 20",
			args: []interface{}{"bob", 3, "%bob", 100},
		},
		{
			goql: z.Select(ID).From(Users).Where(
				z.Exists(z.Select(z.Var(1))),
			),
			sql:  "SELECT `users`.`id` FROM `users` WHERE EXISTS (SELECT $1)",
			args: []interface{}{1},
		},
		{
			goql: z.Select(z.Select(z.Var(1)).As("subquery")),
			sql:  "SELECT (SELECT $1) AS `subquery`",
			args: []interface{}{1},
		},
		{
			goql: z.Select(z.Var(1)).From(
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
			goql: z.Select(
				z.Col("", "id"), z.Col("f", "title"), z.Col("foo", "body").As("content"),
				z.Col("", "count").Add(3),
			).From(
				z.Table("foo").As("f"),
			),
			sql:  "SELECT `id`, `f`.`title`, `foo`.`body` AS `content`, `count` + $1 FROM `foo` AS `f`",
			args: []interface{}{3},
		},
		{
			goql: z.Select(z.Var(1)).From(Users).Joins(Join(Users).On(Users.ID.Eq(3))),
			sql:  "SELECT $1 FROM `users` INNER JOIN `users` ON `users`.`id` = $2",
			args: []interface{}{1, 3},
		},
		{
			goql: z.Select(z.Var(1).As("num")).OrderBy(z.Name("num")),
			sql:  "SELECT $1 AS `num` ORDER BY `num`",
			args: []interface{}{1},
		},
		{
			goql: z.Select(z.Var(1)).OrderBy(Name.Asc(), ID.Desc(), ID),
			sql:  "SELECT $1 ORDER BY `users`.`name`, `users`.`id` DESC, `users`.`id`",
			args: []interface{}{1},
		},
		{
			// Postgres needs type information for placeholders in some cases.
			goql: z.Select(z.VarT(1, "int")),
			sql:  "SELECT $1::int",
			args: []interface{}{1},
		},
		{
			goql: z.Select(z.Null().As("a")).From(Users).Where(z.Null().Eq(z.Null())),
			sql:  "SELECT NULL AS `a` FROM `users` WHERE NULL = NULL",
			args: nil,
		},
		{
			goql: z.Select(Users.Except(Users.ID), z.Null()).From(Users),
			sql:  "SELECT `users`.`name`, NULL FROM `users`",
			args: nil,
		},
		{
			goql: z.Select(z.Raw("selected")).From(Users).Where(
				ID.Eq(1), z.Raw("some-predicate"),
			),
			sql:  "SELECT selected FROM `users` WHERE `users`.`id` = $1 AND some-predicate",
			args: []interface{}{1},
		},
		{
			goql: z.InsertInto(Users).ValuesMap(Values{
				Users.ID:   1,
				Users.Name: "bob",
			}),
			sql:  "INSERT INTO `users` VALUES ($1, $2)",
			args: []interface{}{1, "bob"},
		},
		{
			goql: z.InsertInto(Users).Values(user{1, "bob"}),
			sql:  "INSERT INTO `users` VALUES ($1, $2)",
			args: []interface{}{1, "bob"},
		},
		{
			goql: z.InsertInto(Users, Users.Name).Values(user{1, "bob"}),
			sql:  "INSERT INTO `users` (`name`) VALUES ($1)",
			args: []interface{}{"bob"},
		},
		{
			goql: z.InsertInto(
				Users,
				Users.Except(Users.Name).Columns()...,
			).Values(user{1, "bob"}),
			sql:  "INSERT INTO `users` (`id`) VALUES ($1)",
			args: []interface{}{1},
		},
		{
			goql: z.Update(Users).Set(Values{
				Users.ID:   30,
				Users.Name: "alice",
			}).Where(Users.ID.Eq(5)),
			sql:  "UPDATE `users` SET `id` = $1, `name` = $2 WHERE `users`.`id` = $3",
			args: []interface{}{30, "alice", 5},
		},
		{
			goql: z.Update(Users).Elem(user{45, "john"}),
			sql:  "UPDATE `users` SET `id` = $1, `name` = $2 WHERE `users`.`id` = $3",
			args: []interface{}{45, "john", 45},
		},
		{
			goql: z.Update(Users).Elem(user{45, "john"}, Users.Name),
			sql:  "UPDATE `users` SET `name` = $1 WHERE `users`.`id` = $2",
			args: []interface{}{"john", 45},
		},
		{
			goql: z.DeleteFrom(Users),
			sql:  "DELETE FROM `users`",
			args: nil,
		},
	}

	for i, test := range tests {
		q := z.Query(test.goql)
		if len(q.errs) > 0 {
			for _, err := range q.errs {
				t.Error(err)
			}
		} else {
			if query := q.Query(); query != test.sql {
				t.Errorf("[%d] Query diff\nGOT : %s\nWANT: %s", i, query, test.sql)
			}
			if diff := cmp.Diff(q.Args(), test.args); diff != "" {
				t.Errorf("[%d] Args diff\n%s", i, diff)
			}
		}
	}
}
