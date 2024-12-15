package reflectx

import (
	"database/sql"
	"reflect"
	"time"
	"unsafe"
)

type FieldPtrGetter func(insptr unsafe.Pointer) any

func appendType[T any, M any](tmap map[reflect.Type]FieldPtrGetter, field *Field[M]) {
	tmap[Typeof[T]()] = func(insptr unsafe.Pointer) any { return (*T)(unsafe.Add(insptr, field.Offset)) }
	tmap[Typeof[*T]()] = func(insptr unsafe.Pointer) any { return (**T)(unsafe.Add(insptr, field.Offset)) }
	tmap[Typeof[sql.Null[T]]()] = func(insptr unsafe.Pointer) any { return (*sql.Null[T])(unsafe.Add(insptr, field.Offset)) }
}

func (field *Field[M]) PtrGetter() FieldPtrGetter {
	builtinTypes := map[reflect.Type]FieldPtrGetter{}

	appendType[string](builtinTypes, field)
	appendType[[]byte](builtinTypes, field)

	appendType[int8](builtinTypes, field)
	appendType[int16](builtinTypes, field)
	appendType[int32](builtinTypes, field)
	appendType[int64](builtinTypes, field)

	appendType[uint8](builtinTypes, field)
	appendType[uint16](builtinTypes, field)
	appendType[uint32](builtinTypes, field)
	appendType[uint64](builtinTypes, field)

	appendType[int](builtinTypes, field)
	appendType[uint](builtinTypes, field)

	appendType[float32](builtinTypes, field)
	appendType[float64](builtinTypes, field)

	appendType[bool](builtinTypes, field)

	appendType[time.Time](builtinTypes, field)
	appendType[time.Duration](builtinTypes, field)

	g, ok := builtinTypes[field.Field.Type]
	if ok {
		return g
	}
	return func(insptr unsafe.Pointer) any {
		return reflect.NewAt(field.Field.Type, unsafe.Add(insptr, field.Offset)).Interface()
	}
}
