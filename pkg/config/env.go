package config

import (
	"fmt"
	"reflect"
)

type EnvItem struct {
	Name  string
	Value string
}

func ToEnvFromObject[T any](data T) []EnvItem {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	var items []EnvItem
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		raw := v.Field(i).Interface()

		items = append(items, EnvItem{
			Name:  name,
			Value: fmt.Sprintf("%v", raw),
		})
	}
	return items
}
