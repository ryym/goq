package cllct

import (
	"fmt"
	"reflect"
)

func checkPtrKind(ptr interface{}, kind reflect.Kind) error {
	tp := reflect.TypeOf(ptr)
	if tp.Kind() != reflect.Ptr || tp.Elem().Kind() != kind {
		return fmt.Errorf("required: pointer of %s, got: %s", kind, tp)
	}
	return nil
}

func checkSliceMapPtrKind(ptr interface{}) error {
	tp := reflect.TypeOf(ptr)
	if tp.Kind() == reflect.Ptr {
		elem := tp.Elem()
		if elem.Kind() == reflect.Map && elem.Elem().Kind() == reflect.Slice {
			return nil
		}
	}
	return fmt.Errorf("required: pointer of a map of slices, got: %s", tp)
}
