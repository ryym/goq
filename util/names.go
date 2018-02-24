// Package util is an internal utility.
//
// This package provides a name convertion between
// a column and a field. Currently Goq assumes that
// a field name uses pascal case (e.g. UserID) and
// a column name uses snake case (e.g. user_id).
package util

import "unicode"

// FldToCol converts a field name to column name.
func FldToCol(name string) string {
	if name == "" {
		return name
	}

	fld := []rune(name)
	col := []rune{unicode.ToLower(fld[0])}

	rs := append(fld, 'A') // Add dummy padding.
	for i := 1; i < len(fld); i++ {
		if unicode.IsUpper(rs[i]) {
			if !unicode.IsUpper(rs[i-1]) || !unicode.IsUpper(rs[i+1]) {
				col = append(col, '_')
			}
			col = append(col, unicode.ToLower(rs[i]))
		} else {
			col = append(col, rs[i])
		}
	}

	return string(col)
}

// foo_bar_baz -> FooBarBaz
// NOTE: impossible: `id -> ID`

// ColToFld converts a column name to field name.
// This conversion will be incorrect in some cases.
// Example:
//     FldToCol("UserID") //=> "user_id"
//     ColToFld("user_id") //=> "UserId"
func ColToFld(name string) string {
	if name == "" {
		return name
	}

	col := []rune(name)
	fld := []rune{unicode.ToUpper(col[0])}

	for i := 1; i < len(col); i++ {
		if col[i] == '_' {
			i++
			fld = append(fld, unicode.ToUpper(col[i]))
		} else {
			fld = append(fld, col[i])
		}
	}

	return string(fld)
}
