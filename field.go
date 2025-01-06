package lion

import (
	"fmt"
	"reflect"
	"unsafe"
)

type _FieldPtrGetter func(insptr unsafe.Pointer) any

type Field[M any] struct {
	typeinfo  *TypeInfo[M]
	offset    int64
	ref       *Field[M]
	name      string
	field     reflect.StructField
	meta      *M
	ptrgetter _FieldPtrGetter
	getter    _FieldPtrGetter
	setter    func(insptr unsafe.Pointer, val any)
}

func (field *Field[M]) Typeinfo() *TypeInfo[M] {
	return field.typeinfo
}

// UnsafeUpdate
// fast but unsafe, you must know the field's type
func UnsafeUpdate[T any, M any, V any](insptr *T, field *Field[M], val V) {
	fuptr := unsafe.Add(unsafe.Pointer(insptr), field.offset)
	fptr := (*V)(fuptr)
	*fptr = val
}

// UnsafeFieldPtr
// same as `UnsafeUpdate`
func UnsafeFieldPtr[T any, M any, V any](insptr *T, field *Field[M]) *V {
	return (*V)(unsafe.Add(unsafe.Pointer(insptr), field.offset))
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

func (field *Field[M]) String() string {
	sf := field.StructField()
	return fmt.Sprintf("Field{Offset: %d, ReflectName: %s, StructName: %s, StructType: %s}", field.Offset(), field.Name(), sf.Name, sf.Type)
}

func FieldOf[T any, M any](ptr any) *Field[M] {
	return TypeInfoOf[T, M]().FieldByPtr(ptr)
}
