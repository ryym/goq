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
		, t.id
		, t.name
	from
		prefectures p
	inner join cities c
		on c.prefecture_id = p.id
	inner join towns t
		on t.city_id = c.id
	order by
		c.name
	limit 50
	`)
	chk(err)
	defer rows.Close()

	cols := []Column{
		{"Pref", "ID"},
		{"Pref", "Name"},
		{"City", "ID"},
		{"City", "Name"},
		{"Town", "ID"},
		{"Town", "Name"},
	}

	prefs := []Pref{}
	// cities := []City{}
	// towns := []Town{}
	cities := map[int][]City{}
	towns := map[int][]Town{}

	colls := []Collector{
		NewSliceCollector(cols, &prefs, reflect.TypeOf(Pref{})),
		// NewSliceCollector(cols, &cities, reflect.TypeOf(City{})),
		// NewSliceCollector(cols, &towns, reflect.TypeOf(Town{})),

		NewMapCollector(cols, &cities, reflect.TypeOf(City{}), Column{"Pref", "ID"}),
		NewMapCollector(cols, &towns, reflect.TypeOf(Town{}), Column{"City", "ID"}),
	}

	ptrs := make([]interface{}, len(cols))
	for rows.Next() {
		for _, cl := range colls {
			cl.Next(ptrs)
		}
		rows.Scan(ptrs...)
		for _, cl := range colls {
			cl.AfterScan(ptrs)
		}
	}

	fmt.Println(prefs)
	fmt.Println(cities)
	fmt.Println(towns)

	err = rows.Err()
	chk(err)
}

type SliceCollector struct {
	tp      reflect.Type
	colToFl map[int]int

	slice reflect.Value
	item  reflect.Value
}

type Collector interface {
	Next(ptrs []interface{})
	AfterScan(ptrs []interface{})
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

func (sc *SliceCollector) AfterScan(ptrs []interface{}) {
	sc.slice.Set(reflect.Append(sc.slice, sc.item.Elem()))
}

type MapCollector struct {
	tp      reflect.Type
	colToFl map[int]int
	keyIdx  int
	mp      reflect.Value
	item    reflect.Value
}

func NewMapCollector(cols []Column, mp interface{}, itemType reflect.Type, keyCol Column) *MapCollector {
	names := make([]string, itemType.NumField())
	for i := 0; i < itemType.NumField(); i++ {
		names[i] = itemType.Field(i).Name
	}

	keyIdx := -1
	colToFl := map[int]int{}
	for iC, c := range cols {
		if c.StructName == itemType.Name() {
			for iF, name := range names {
				if name == c.FieldName {
					colToFl[iC] = iF
				}
			}
		}
		if c == keyCol {
			keyIdx = iC
		}
	}

	if keyIdx == -1 {
		panic("[NewMapCollector]: key is not in columns")
	}

	return &MapCollector{
		tp:      itemType,
		colToFl: colToFl,
		keyIdx:  keyIdx,
		mp:      reflect.ValueOf(mp).Elem(),
	}
}

func (sc *MapCollector) Next(ptrs []interface{}) {
	item := reflect.New(sc.tp).Elem()
	sc.item = item.Addr()

	for c, f := range sc.colToFl {
		ptrs[c] = item.Field(f).Addr().Interface()
	}
}

func (sc *MapCollector) AfterScan(ptrs []interface{}) {
	key := reflect.ValueOf(ptrs[sc.keyIdx]).Elem()

	sl := sc.mp.MapIndex(key)
	if !sl.IsValid() {
		sl = reflect.MakeSlice(reflect.SliceOf(sc.tp), 0, 0)
	}
	sl = reflect.Append(sl, sc.item.Elem())
	sc.mp.SetMapIndex(key, sl)
}
