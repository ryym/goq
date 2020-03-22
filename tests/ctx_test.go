package tests

import "strconv"

const (
	DB_POSTGRES = "postgres"
	DB_MYSQL    = "mysql"
)

type testCtx struct {
	dbName string
}

func (c *testCtx) rawStr(v string) interface{} {
	switch c.dbName {
	case DB_POSTGRES:
		return v
	default:
		return []uint8(v)
	}
}

func (c *testCtx) rawInt(n int) interface{} {
	switch c.dbName {
	case DB_MYSQL:
		return []uint8(strconv.Itoa(n))
	default:
		return int64(n)
	}
}
