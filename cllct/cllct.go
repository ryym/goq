/*
Package cllct provides various collectors.
A collector collects *sql.Rows fetched from DB
into a common data structure like a slice, map, etc.

Available collectors:

	Table.ToElem
	Table.ToSlice
	Table.ToUniqSlice
	Table.ToMap
	Table.ToSliceMap
	Table.ToUniqSliceMap
	ToElem
	ToSlice
	ToMap
	ToSliceMap
	ToRowMap
	ToRowMapSlice

These collector methods can be used from the query builder.
'Table' in the above list means a table helper
(see [TODO: link] for custom query builder).

For example, the following code collects users table records
into a 'users' slice and other values into another slice.

	b := NewBuilder(db.Dialect())
	q := b.Select(
		b.Users.All(),
		b.Coalesce(b.Users.NickName, b.Users.Name).As("nick_name"),
		b.Func("MONTH", b.Users.BirthDay).As("month"),
	).From(b.Users)

	db.Query(q).Collect(
		b.Users.ToSlice(&users),
		b.ToSlice(&others),
	)

The collectors defined in a table helper are called model collector.
They collect rows into a model or models.
The model collector structs are named 'ModelXxxCollector'.
For example, 'Table.ToSlice' returns a ModelSliceCollector.

The collectors defined in a query builder are called generic collector.
They collect rows into a generic structure such as a slice of non-model struct.
The generic collector structs are named 'XxxCollector'.
For example, 'ToSlice' returns a SliceCollector.

Furthermore, collectors are classified into two types: list collector or single collector.
A list collector collects rows into a slice or a map of slices.
A single collector collects a first row only into a struct or map.
You need to pass them to Collect and First methods, respectively.

	db.Query(q).Collect(z.Users.ToSlice(&users))
	db.Query(q).First(z.Users.ToElem(&user))

Note that we often use 'z' as a variable name of Goq query builder in example code.
This name has no special meanings. We use it just because
this character rarely duplicates with other common variable names
and is easy to identify.
*/
package cllct

import "github.com/ryym/goq/goql"

type initConf struct {
	Selects  []goql.Selection
	ColNames []string
	takens   map[int]bool
}

// NewInitConf creates a initConf for collectors.
// This is used internally.
func NewInitConf(selects []goql.Selection, colNames []string) *initConf {
	return &initConf{selects, colNames, map[int]bool{}}
}

func (c *initConf) take(colIdx int) bool {
	ok := c.takens[colIdx]
	if !ok {
		c.takens[colIdx] = true
	}
	return !ok
}

func (c *initConf) canTake(colIdx int) bool {
	return !c.takens[colIdx]
}

// Collector defines methods to collect query results.
type Collector interface {
	next(ptrs []interface{})
	init(conf *initConf) (collectable bool, err error)
	afterScan(ptrs []interface{})
	afterinit(conf *initConf) error
}

// ListCollector interface represents a collector which
// collects rows into a collection data.
type ListCollector interface {
	Collector
	ImplListCollector()
}

// SingleCollector interface represents a collector which
// scans a first row.
type SingleCollector interface {
	Collector
	ImplSingleCollector()
}

type tableInfo struct {
	structName string
	tableAlias string
}
