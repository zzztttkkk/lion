package lion

import (
	"fmt"
	"reflect"
	"unsafe"
)

type FieldPtrGetter func(insptr unsafe.Pointer) any

type Field struct {
	typeinfo *TypeInfo
	offset   int64
	ref      *Field
	field    reflect.StructField
	metas    map[reflect.Type]any
	tags     map[string]*Tag

	ptrgetter FieldPtrGetter
	setter    func(insptr unsafe.Pointer, val any)
}

func (field *Field) Typeinfo() *TypeInfo {
	return field.typeinfo
}

func (field *Field) Offset() int64 {
	return field.offset
}

func (field *Field) StructField() *reflect.StructField {
	if field.ref == nil {
		return &field.field
	}
	return &field.ref.field
}

func (filed *Field) getMetainfo(metatype reflect.Type) any {
	if filed.ref == nil {
		return filed.metas[metatype]
	}
	return filed.ref.getMetainfo(metatype)
}

func (filed *Field) updateMetainfo(metatype reflect.Type, meta any) {
	if filed.ref == nil {
		if filed.metas == nil {
			filed.metas = map[reflect.Type]any{}
		}
		filed.metas[metatype] = meta
		return
	}
	filed.ref.updateMetainfo(metatype, meta)
}

func UpdateMetainfo[T any, M any](fptr any, meta *M) {
	FieldOf[T](fptr).updateMetainfo(Typeof[M](), meta)
}

func MetainfoOf[T any, M any](fptr any) *M {
	val := FieldOf[T](fptr).getMetainfo(Typeof[M]())
	if val == nil {
		return nil
	}
	return val.(*M)
}

func (field *Field) String() string {
	sf := field.StructField()
	return fmt.Sprintf("Field{Offset: %d, StructName: %s, StructType: %s}", field.Offset(), sf.Name, sf.Type)
}

func FieldOf[T any](ptr any) *Field {
	return TypeInfoOf[T]().FieldByPtr(ptr)
}
