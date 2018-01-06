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

	cols := []Column{
		{"User", "Name"},
	}

	users := []User{}
	coll := NewCollector(cols, &users)
	ptrs := make([]interface{}, len(cols))
	for rows.Next() {
		coll.Next(ptrs)
		rows.Scan(ptrs...)
		coll.AfterScan()
	}

	fmt.Println("result", users)

	err = rows.Err()
	chk(err)
}

type UserSliceCollector struct {
	ut      reflect.Type
	colToFl map[int]int

	users reflect.Value
	user  reflect.Value
}

func NewCollector(cols []Column, users interface{}) *UserSliceCollector {
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
		users:   reflect.ValueOf(users).Elem(),
	}
}

func (uc *UserSliceCollector) Next(ptrs []interface{}) {
	user := reflect.New(uc.ut).Elem()
	uc.user = user.Addr()

	for c, f := range uc.colToFl {
		ptrs[c] = user.Field(f).Addr().Interface()
	}
}

// XXX: AfterScan 的な処理はどうしても必要そう
// 実際にはコレクタは複数いるから、rows のループ内で
// コレクタのループを scan 前後に計2回も実行しないといけない...
// モデルのポインタのスライスでよければ、`Next`内だけで作れる。
// しかしモデルの実体のスライスは無理っぽい。
// フィールドのポインタに値がセットされる前に`Elem`や`Interface`を
// 呼ぶと空のモデルになっちゃう。
// After ではなく Next の前に1つ前のループの後処理を行えばコレクタのループは
// 1回で済むかも。
func (uc *UserSliceCollector) AfterScan() {
	uc.users.Set(reflect.Append(uc.users, uc.user.Elem()))
	// fmt.Println(uc.users)
}
