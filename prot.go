package main

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

// とりあえず結果マッピングの構想を試したい。
// - select のカラムリスト
// - collector
// - テーブル名とカラム名を抽出できるモデル
// - (snake <-> camel)

type User struct {
	Name string
}

type Column struct {
	StructName string
	FieldName  string

	// クエリ作成のためにこれらも必要
	// TableName string
	// ColumnName string
}

type UsersSchema struct {
	Name Column
}

var Users = UsersSchema{
	Name: Column{"users", "name"},
}

func main() {
	prot()
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func prot() {
	db, err := sql.Open("sqlite3", "prot.db")
	chk(err)
	defer db.Close()

	rows, err := db.Query("select * from users")
	chk(err)
	defer rows.Close()

	// for rows.Next() {
	// 	user := User{}
	// 	err = rows.Scan(&user.Name)
	// 	chk(err)
	// 	// fmt.Println(user)
	// }

	cols := []Column{
		{"User", "Name"},
	}
	coll := NewCollector(cols)
	for rows.Next() {
		// coll.Next(rows)
		ptrs := make([]interface{}, len(cols))
		coll.Next(ptrs)
		rows.Scan(ptrs...)
		fmt.Println("user", coll.user)
	}

	err = rows.Err()
	chk(err)
}

type UserSliceCollector struct {
	ut      reflect.Type
	colToFl map[int]int
	user    *User // XXX
}

func NewCollector(cols []Column) *UserSliceCollector {
	ut := reflect.TypeOf(User{})

	names := make([]string, ut.NumField())
	for i := 0; i < ut.NumField(); i++ {
		names[i] = ut.Field(i).Name
	}

	colToFl := map[int]int{}
	for iC, c := range cols {
		if c.StructName == "User" {
			for iF, name := range names {
				if name == c.FieldName {
					colToFl[iC] = iF
				}
			}
		}
	}

	return &UserSliceCollector{
		ut:      ut,
		colToFl: colToFl,
	}
}

// func (uc *UserSliceCollector) Next(rows *sql.Rows) {
func (uc *UserSliceCollector) Next(ptrs []interface{}) {
	user := reflect.New(uc.ut).Elem()
	uc.user = user.Addr().Interface().(*User)

	// rows.Scan(user.Field(0).Addr().Interface())
	// uc.users = append(uc.users, user.Interface())
	// fmt.Println(user.Interface())

	for c, f := range uc.colToFl {
		ptrs[c] = user.Field(f).Addr().Interface()
	}
}
