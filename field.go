package reflectx

import (
	"reflect"
)

type Field struct {
	Name   string
	Offset int64
	Field  reflect.StructField
	metas  map[reflect.Type]any
}

func FieldOf[T any](ptr any) *Field {
	return TypeInfoOf[T]().FieldByPtr(ptr)
}
