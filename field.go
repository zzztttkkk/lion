package lion

import (
	"fmt"
	"reflect"
	"unsafe"
)

type _FieldPtrGetter func(insptr unsafe.Pointer) any

// Field
// represents the field information of a struct.
type Field struct {
	typeinfo *TypeInfo
	offset   int64
	ref      *Field
	field    reflect.StructField
	metas    map[reflect.Type]any
	tags     map[string]*Tag

	ptrgetter _FieldPtrGetter
	setter    func(insptr unsafe.Pointer, val any)
}

// TypeInfo
// returns the type information of the field's struct.
func (field *Field) TypeInfo() *TypeInfo {
	return field.typeinfo
}

// Offset
// returns the offset of the field in the struct.
func (field *Field) Offset() int64 {
	return field.offset
}

// StructField
// returns the *reflect.StructField of the field.
func (field *Field) StructField() *reflect.StructField {
	if field.ref == nil {
		return &field.field
	}
	return &field.ref.field
}

func (field *Field) String() string {
	sf := field.StructField()
	return fmt.Sprintf("Field{Offset: %d, StructName: %s, StructType: %s}", field.Offset(), sf.Name, sf.Type)
}

// FieldOf
// returns the ptr's field information of the struct type T.
func FieldOf[T any](ptr any) *Field {
	return TypeInfoOf[T]().FieldByPtr(ptr)
}
