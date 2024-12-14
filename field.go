package reflectx

import (
	"reflect"
	"unsafe"
)

type Field[M any] struct {
	Name   string
	Offset int64
	Field  reflect.StructField
	Meta   *M
}

func (field *Field[M]) GetPtrValue(insptr unsafe.Pointer) reflect.Value {
	return reflect.NewAt(field.Field.Type, unsafe.Add(insptr, field.Offset))
}

func (field *Field[M]) Set(insptr unsafe.Pointer, val reflect.Value) {
	field.GetPtrValue(insptr).Elem().Set(val)
}

func (field *Field[M]) SetAny(insptr unsafe.Pointer, val any) {
	field.Set(insptr, reflect.ValueOf(val))
}

func FieldOf[T any, M any](ptr any) *Field[M] {
	return TypeInfoOf[T, M]().FieldByPtr(ptr)
}
