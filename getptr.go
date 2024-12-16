package reflectx

import (
	"database/sql"
	"reflect"
	"time"
	"unsafe"
)

var (
	builtinTypePtrGetters = map[reflect.Type]func(insptr unsafe.Pointer, offset int64) any{}
)

func appendType[T any]() {
	builtinTypePtrGetters[Typeof[T]()] = func(insptr unsafe.Pointer, offset int64) any { return (*T)(unsafe.Add(insptr, offset)) }
	builtinTypePtrGetters[Typeof[*T]()] = func(insptr unsafe.Pointer, offset int64) any { return (**T)(unsafe.Add(insptr, offset)) }
	builtinTypePtrGetters[Typeof[sql.Null[T]]()] = func(insptr unsafe.Pointer, offset int64) any { return (*sql.Null[T])(unsafe.Add(insptr, offset)) }
}

func init() {
	appendType[string]()
	appendType[[]byte]()

	appendType[int8]()
	appendType[int16]()
	appendType[int32]()
	appendType[int64]()

	appendType[uint8]()
	appendType[uint16]()
	appendType[uint32]()
	appendType[uint64]()

	appendType[int]()
	appendType[uint]()

	appendType[float32]()
	appendType[float64]()

	appendType[bool]()

	appendType[time.Time]()
	appendType[time.Duration]()
}

func (field *Field[M]) PtrGetter() _FieldPtrGetter {
	if field.ptrgetter == nil {
		getter, ok := builtinTypePtrGetters[field.Field.Type]
		if ok {
			field.ptrgetter = func(insptr unsafe.Pointer) any { return getter(insptr, field.Offset) }
		} else {
			field.ptrgetter = func(insptr unsafe.Pointer) any {
				return reflect.NewAt(field.Field.Type, unsafe.Add(insptr, field.Offset)).Interface()
			}
		}
	}
	return field.ptrgetter
}

func (field *Field[M]) PtrOf(insptr unsafe.Pointer) any { return field.PtrGetter()(insptr) }
