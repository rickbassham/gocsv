package gocsv

import (
	"reflect"
	"strconv"
	"strings"
)

func parseTag(tag string) (string, []string) {
	split := strings.Split(tag, ",")

	if len(split) > 0 {
		return split[0], split[1:]
	}

	return "", []string{}
}

func omitEmpty(tagOptions []string) bool {
	for _, option := range tagOptions {
		if option == "omitempty" {
			return true
		}
	}

	return false
}

func base(tag reflect.StructTag) (int, error) {
	base := int64(10)

	if baseStr, ok := tag.Lookup("base"); ok {
		var err error
		base, err = strconv.ParseInt(baseStr, 10, 32)
		if err != nil {
			return int(base), ErrInvalidIntBase
		}
	}

	return int(base), nil
}
