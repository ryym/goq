package main

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

// 複数テーブルへの結果マッピング
// - 1:1
// - 1:N
// - N:N
// ネストした結合のマッピング

type Pref struct {
	ID   int
	Name string
}
type City struct {
	ID     int
	Name   string
	PrefID int
}
type Town struct {
	ID     int
	Name   string
	CityID int
}

var Prefs = struct {
	ID   Column
	Name Column
}{
	Column{"prefectures", "id"},
	Column{"prefectures", "name"},
}
var Cities = struct {
	ID     Column
	Name   Column
	PrefID Column
}{
	Column{"cities", "id"},
	Column{"cities", "name"},
	Column{"cities", "prefecture_id"},
}
var Towns = struct {
	ID     Column
	Name   Column
	CityID Column
}{
	Column{"towns", "id"},
	Column{"towns", "name"},
	Column{"towns", "city_id"},
}

func prot2() {
	db, err := sql.Open("sqlite3", "prot.db")
	chk(err)
	defer db.Close()

	// rows, err := db.Query("select id,name from prefectures limit 10")
	rows, err := db.Query(`
	select
		p.id
		, p.name
		, c.id
		, c.name
	from
		prefectures p
	inner join cities c
		on c.prefecture_id = p.id
	limit 10
	`)
	chk(err)
	defer rows.Close()

	cols := []Column{
		{"Pref", "ID"},
		{"Pref", "Name"},
		{"City", "ID"},
		{"City", "Name"},
	}

	prefs := []Pref{}
	cities := []City{}
	prefC := NewSliceCollector(cols, &prefs, reflect.TypeOf(Pref{}))
	cityC := NewSliceCollector(cols, &cities, reflect.TypeOf(City{}))
	ptrs := make([]interface{}, len(cols))
	for rows.Next() {
		prefC.Next(ptrs)
		cityC.Next(ptrs)
		rows.Scan(ptrs...)
		prefC.AfterScan()
		cityC.AfterScan()
	}

	fmt.Println(prefs)
	fmt.Println(cities)

	// err = rows.Err()
	// chk(err)
}

type SliceCollector struct {
	tp      reflect.Type
	colToFl map[int]int

	slice reflect.Value
	item  reflect.Value
}

func NewSliceCollector(cols []Column, slice interface{}, itemType reflect.Type) *SliceCollector {
	names := make([]string, itemType.NumField())
	for i := 0; i < itemType.NumField(); i++ {
		names[i] = itemType.Field(i).Name
	}

	colToFl := map[int]int{}
	for iC, c := range cols {
		if c.StructName == itemType.Name() {
			for iF, name := range names {
				if name == c.FieldName {
					colToFl[iC] = iF
				}
			}
		}
	}

	return &SliceCollector{
		tp:      itemType,
		colToFl: colToFl,
		slice:   reflect.ValueOf(slice).Elem(),
	}
}

func (sc *SliceCollector) Next(ptrs []interface{}) {
	item := reflect.New(sc.tp).Elem()
	sc.item = item.Addr()

	for c, f := range sc.colToFl {
		ptrs[c] = item.Field(f).Addr().Interface()
	}
}

func (sc *SliceCollector) AfterScan() {
	sc.slice.Set(reflect.Append(sc.slice, sc.item.Elem()))
}
