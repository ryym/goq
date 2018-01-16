package util

import "unicode"

func FldToCol(name string) string {
	if name == "" {
		return name
	}

	fld := []rune(name)
	col := []rune{unicode.ToLower(fld[0])}

	rs := append(fld, 'A') // Add dummy padding.
	for i := 1; i < len(fld); i++ {
		if unicode.IsUpper(rs[i]) {
			// 前後が両方 upper なら`_`は不要。
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
