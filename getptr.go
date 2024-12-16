package reflectx

import (
	"database/sql"
	"reflect"
	"time"
	"unsafe"
)

var (
	preparedTypePtrGetters = map[reflect.Type]func(insptr unsafe.Pointer, offset int64) any{}
)

func appendType[T any]() {
	preparedTypePtrGetters[Typeof[T]()] = func(insptr unsafe.Pointer, offset int64) any { return (*T)(unsafe.Add(insptr, offset)) }
	preparedTypePtrGetters[Typeof[*T]()] = func(insptr unsafe.Pointer, offset int64) any { return (**T)(unsafe.Add(insptr, offset)) }

	preparedTypePtrGetters[Typeof[[]T]()] = func(insptr unsafe.Pointer, offset int64) any { return (*[]T)(unsafe.Add(insptr, offset)) }
	preparedTypePtrGetters[Typeof[[]*T]()] = func(insptr unsafe.Pointer, offset int64) any { return (*[]*T)(unsafe.Add(insptr, offset)) }

	preparedTypePtrGetters[Typeof[sql.Null[T]]()] = func(insptr unsafe.Pointer, offset int64) any { return (*sql.Null[T])(unsafe.Add(insptr, offset)) }
}

func AddType[T any]() {
	appendType[T]()
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
		sf := field.StructField()
		getter, ok := preparedTypePtrGetters[sf.Type]
		if ok {
			field.ptrgetter = func(insptr unsafe.Pointer) any { return getter(insptr, field.offset) }
		} else {
			field.ptrgetter = func(insptr unsafe.Pointer) any {
				return reflect.NewAt(field.field.Type, unsafe.Add(insptr, field.offset)).Interface()
			}
		}
	}
	return field.ptrgetter
}

func (field *Field[M]) PtrOf(insptr unsafe.Pointer) any { return field.PtrGetter()(insptr) }
