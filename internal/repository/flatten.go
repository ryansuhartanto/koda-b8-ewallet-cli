package repository

import (
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
)

func StrictFlattenArgs(obj any) pgx.StrictNamedArgs {
	out := make(pgx.StrictNamedArgs)
	flattenArgs(out, reflect.ValueOf(obj))
	return out
}

func FlattenArgs(obj any) pgx.NamedArgs {
	out := make(pgx.NamedArgs)
	flattenArgs(out, reflect.ValueOf(obj))
	return out
}

func flattenArgs(out map[string]any, v reflect.Value) {
	t := v.Type()
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		fv := v.Field(i)
		if f.Anonymous {
			if fv.Kind() == reflect.Pointer {
				if fv.IsNil() {
					continue
				}
				fv = fv.Elem()
			}
			if fv.Kind() == reflect.Struct {
				flattenArgs(out, fv)
				continue
			}
		}
		key := f.Name
		if tag, ok := f.Tag.Lookup("db"); ok {
			tag, _, _ = strings.Cut(tag, ",")
			if tag == "-" {
				continue
			}
			key = tag
		}
		out[key] = fv.Interface()
	}
}
