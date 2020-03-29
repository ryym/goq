package goq

type initConf struct {
	Selects  []Selection
	ColNames []string
	takens   map[int]bool
}

// NewInitConf creates a initConf for collectors.
// This is used internally.
func NewInitConf(selects []Selection, colNames []string) *initConf {
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
