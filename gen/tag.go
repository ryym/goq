package gen

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	TYPE_STR  = "string"
	TYPE_BOOL = "bool"
)

type Type string

func ParseTag(tag string, wants map[string]Type) (map[string]interface{}, error) {
	tagRgx := regexp.MustCompile(";")
	matches := tagRgx.Split(tag, -1)

	ret := make(map[string]interface{}, len(matches))
	for _, m := range matches {
		if len(m) == 0 {
			continue
		}

		kv := strings.Split(m, ":")
		name := kv[0]

		tp, wanted := wants[name]
		if !wanted {
			return nil, fmt.Errorf("Unkown tag attribute: %s", name)
		}

		if len(kv) == 1 {
			if tp != TYPE_BOOL {
				return nil, fmt.Errorf("'%s' must be bool", name)
			}
			ret[name] = true
		} else {
			if tp != TYPE_STR {
				return nil, fmt.Errorf("'%s' must have a string value", name)
			}
			ret[name] = kv[1]
		}
	}

	return ret, nil
}

type ModelTag struct {
	IsPK    bool
	ColName string
	NotCol  bool
}

func ParseModelTag(tagStr string) (ModelTag, error) {
	attrs, err := ParseTag(tagStr, map[string]Type{
		"pk":   TYPE_BOOL,
		"name": TYPE_STR,
		"-":    TYPE_BOOL,
	})
	if err != nil {
		return ModelTag{}, err
	}

	if _, ok := attrs["-"]; ok && len(attrs) > 1 {
		return ModelTag{}, errors.New(`"-" cannot be used with other attributes`)
	}

	return makeModelTag(attrs), nil
}

func makeModelTag(attrs map[string]interface{}) ModelTag {
	tag := ModelTag{}
	for key, val := range attrs {
		switch key {
		case "pk":
			tag.IsPK = val.(bool)
		case "name":
			tag.ColName = val.(string)
		case "-":
			tag.NotCol = val.(bool)
		}
	}
	return tag
}

type TableTag struct {
	HelperName string
}

func ParseTableTag(tagStr string) (TableTag, error) {
	attrs, err := ParseTag(tagStr, map[string]Type{
		"helper": TYPE_STR,
	})
	if err != nil {
		return TableTag{}, err
	}

	return makeTableTag(attrs), nil
}

func makeTableTag(attrs map[string]interface{}) TableTag {
	tag := TableTag{}
	for key, val := range attrs {
		switch key {
		case "helper":
			tag.HelperName = val.(string)
		}
	}
	return tag
}
