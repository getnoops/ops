package config

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	titler   = cases.Title(language.Und)
	uuidType = reflect.TypeOf(uuid.UUID{})
)

// EnvSet represents a set of environment variables.
type EnvSet struct {
	// Key is the environment variable key.
	Key string
	// Value is the environment variable value.
	Value string
}

func EnvMarshal(name string, obj interface{}) ([]EnvSet, error) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	// handle special types
	if t == uuidType {
		str := fmt.Sprintf("%v", v.Interface())
		return []EnvSet{{
			Key:   name,
			Value: str,
		}}, nil
	}

	switch t.Kind() {
	case reflect.Ptr:
		return EnvMarshal(name, v.Elem().Interface())

	case reflect.Struct:
		es := []EnvSet{}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			vf := v.Field(i)

			if f.Type.Kind() == reflect.Ptr && vf.IsNil() {
				continue
			}

			p := name
			if p != "" {
				p += "_"
			}
			p += f.Name

			es2, err := EnvMarshal(p, vf.Interface())
			if err != nil {
				return nil, err
			}
			es = append(es, es2...)
		}
		return es, nil

	case reflect.Slice, reflect.Array:
		es := []EnvSet{}
		for i := 0; i < v.Len(); i++ {
			p := name
			if p != "" {
				p += "_"
			}
			p += fmt.Sprintf("%d", i)

			es2, err := EnvMarshal(p, v.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			es = append(es, es2...)
		}
		return es, nil

	case reflect.Map:
		es := []EnvSet{}
		for _, k := range v.MapKeys() {
			key := fmt.Sprintf("%v", k.Interface())
			key = titler.String(key)

			p := name
			if p != "" {
				p += "_"
			}
			p += key

			es2, err := EnvMarshal(p, v.MapIndex(k).Interface())
			if err != nil {
				return nil, err
			}
			es = append(es, es2...)
		}
		return es, nil

	default:
		str := fmt.Sprintf("%v", v.Interface())
		return []EnvSet{{
			Key:   name,
			Value: str,
		}}, nil
	}
}
