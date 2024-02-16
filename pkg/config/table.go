package config

import (
	"fmt"
	"reflect"
)

type TableData struct {
	Headers []string
	Rows    [][]string
}

type TableColumn struct {
	Name     string
	Resolver func(interface{}) string
}

func GetTableColumns(t reflect.Type) []TableColumn {
	var columns []TableColumn
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name

		columns = append(columns, TableColumn{
			Name: name,
			Resolver: func(item interface{}) string {
				val := reflect.ValueOf(item).FieldByName(name).Interface()
				return fmt.Sprintf("%v", val)
			},
		})
	}
	return columns
}

func ToTableFromList[T any](data []T) *TableData {
	var a T
	t := reflect.TypeOf(a)

	columns := GetTableColumns(t)
	headers := make([]string, len(columns))
	for i, column := range columns {
		headers[i] = column.Name
	}

	rows := make([][]string, len(data))
	for i, item := range data {
		var rowValues []string
		for _, column := range columns {
			rowValues = append(rowValues, column.Resolver(item))
		}
		rows[i] = rowValues
	}

	return &TableData{
		Headers: headers,
		Rows:    rows,
	}
}

func ToTableFromObject[T any](data T) *TableData {
	return ToTableFromList[T]([]T{data})
}
