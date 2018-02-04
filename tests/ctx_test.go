package tests

const (
	DB_POSTGRES = "postgres"
	DB_SQLITE3  = "sqlite3"
)

type testCtx struct {
	dbName string
}

func (c *testCtx) rawStr(v string) interface{} {
	switch c.dbName {
	case DB_SQLITE3:
		return []uint8(v)
	default:
		return v
	}
}
