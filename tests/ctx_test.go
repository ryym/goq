package tests

const (
	DB_POSTGRES = "postgres"
)

type testCtx struct {
	dbName string
}

func (c *testCtx) rawStr(v string) interface{} {
	switch c.dbName {
	default:
		return v
	}
}
