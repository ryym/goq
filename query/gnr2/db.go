package gnr2

import (
	"database/sql"
	"fmt"
)

// とりあえずの実装

type DB struct {
	db *sql.DB
}

func (d *DB) Query(query QueryStmt) *Collectable {
	return &Collectable{d.db, query}
}

type Collectable struct {
	db    *sql.DB
	query QueryStmt
}

func (c *Collectable) Collect(collectors ...Collector) error {
	qr := c.query.Query()
	fmt.Println("[LOG]", qr)
	rows, err := c.db.Query(qr.Query, qr.Args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	selects := c.query.GetSelects()
	colNames, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(colNames) != len(selects) {
		return fmt.Errorf("colNames: %d, selects: %d", len(colNames), len(selects))
	}

	clls := make([]Collector, 0, len(collectors))
	for _, cl := range collectors {
		if cl.Init(selects, colNames) {
			clls = append(clls, cl)
		}
	}

	ptrs := make([]interface{}, len(colNames))
	for rows.Next() {
		for _, cl := range clls {
			cl.Next(ptrs)
		}
		rows.Scan(ptrs...)
		for _, cl := range clls {
			cl.AfterScan(ptrs)
		}
	}

	return rows.Err()
}
