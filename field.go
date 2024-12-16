package reflectx

import (
	"reflect"
	"unsafe"
)

type _FieldPtrGetter func(insptr unsafe.Pointer) any

type Field[M any] struct {
	Name   string
	Offset int64
	Field  reflect.StructField
	Meta   *M

	ptrgetter _FieldPtrGetter
}

func (field *Field[M]) PtrValueOf(insptr unsafe.Pointer) reflect.Value {
	return reflect.ValueOf(field.PtrOf(insptr))
}

func (field *Field[M]) Set(insptr unsafe.Pointer, val any) {
	field.PtrValueOf(insptr).Elem().Set(reflect.ValueOf(val))
}

func FieldOf[T any, M any](ptr any) *Field[M] {
	return TypeInfoOf[T, M]().FieldByPtr(ptr)
}
