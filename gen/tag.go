package gen

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	tagTypeStr  = "string"
	tagTypeBool = "bool"
)

type tagType string

func parseTag(tag string, wants map[string]tagType) (map[string]interface{}, error) {
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
			if tp != tagTypeBool {
				return nil, fmt.Errorf("'%s' must be bool", name)
			}
			ret[name] = true
		} else {
			if tp != tagTypeStr {
				return nil, fmt.Errorf("'%s' must have a string value", name)
			}
			ret[name] = kv[1]
		}
	}

	return ret, nil
}

type columnTag struct {
	IsPK    bool
	ColName string
	NotCol  bool
}

func parseColumnTag(tagStr string) (columnTag, error) {
	attrs, err := parseTag(tagStr, map[string]tagType{
		"pk":   tagTypeBool,
		"name": tagTypeStr,
		"-":    tagTypeBool,
	})
	if err != nil {
		return columnTag{}, err
	}

	if _, ok := attrs["-"]; ok && len(attrs) > 1 {
		return columnTag{}, errors.New(`"-" cannot be used with other attributes`)
	}

	return makeModelTag(attrs), nil
}

func makeModelTag(attrs map[string]interface{}) columnTag {
	tag := columnTag{}
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

type tableTag struct {
	HelperName string
}

func parseTableTag(tagStr string) (tableTag, error) {
	attrs, err := parseTag(tagStr, map[string]tagType{
		"helper": tagTypeStr,
	})
	if err != nil {
		return tableTag{}, err
	}

	return makeTableTag(attrs), nil
}

func makeTableTag(attrs map[string]interface{}) tableTag {
	tag := tableTag{}
	for key, val := range attrs {
		switch key {
		case "helper":
			tag.HelperName = val.(string)
		}
	}
	return tag
}
