package reflectx

import (
	"reflect"
	"unsafe"
)

type _FieldPtrGetter func(insptr unsafe.Pointer) any

type Field[M any] struct {
	offset    int64
	ref       *Field[M]
	name      string
	field     reflect.StructField
	meta      *M
	ptrgetter _FieldPtrGetter
}

func (field *Field[M]) GetPtrValueOfInstance(insptr unsafe.Pointer) reflect.Value {
	return reflect.ValueOf(field.PtrOf(insptr))
}

func (field *Field[M]) ChangeInstance(insptr unsafe.Pointer, val any) {
	field.GetPtrValueOfInstance(insptr).Elem().Set(reflect.ValueOf(val))
}

func (field *Field[M]) Offset() int64 {
	return field.offset
}

func (field *Field[M]) Name() string {
	if field.ref == nil {
		return field.name
	}
	return field.ref.name
}

func (field *Field[M]) SetName(name string) {
	if field.ref == nil {
		field.name = name
		return
	}
	field.ref.name = name
}

func (field *Field[M]) StructField() *reflect.StructField {
	if field.ref == nil {
		return &field.field
	}
	return &field.ref.field
}

func (filed *Field[M]) Metainfo() *M {
	if filed.ref == nil {
		return filed.meta
	}
	return filed.ref.meta
}

func (filed *Field[M]) UpdateMetainfo(m *M) {
	if filed.ref == nil {
		filed.meta = m
		return
	}
	filed.ref.meta = m
}

func FieldOf[T any, M any](ptr any) *Field[M] {
	return TypeInfoOf[T, M]().FieldByPtr(ptr)
}
