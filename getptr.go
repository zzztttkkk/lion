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

func AppendType[T any]() {
	preparedTypePtrGetters[Typeof[T]()] = func(insptr unsafe.Pointer, offset int64) any { return (*T)(unsafe.Add(insptr, offset)) }
	preparedTypePtrGetters[Typeof[*T]()] = func(insptr unsafe.Pointer, offset int64) any { return (**T)(unsafe.Add(insptr, offset)) }

	preparedTypePtrGetters[Typeof[[]T]()] = func(insptr unsafe.Pointer, offset int64) any { return (*[]T)(unsafe.Add(insptr, offset)) }
	preparedTypePtrGetters[Typeof[[]*T]()] = func(insptr unsafe.Pointer, offset int64) any { return (*[]*T)(unsafe.Add(insptr, offset)) }

	preparedTypePtrGetters[Typeof[sql.Null[T]]()] = func(insptr unsafe.Pointer, offset int64) any { return (*sql.Null[T])(unsafe.Add(insptr, offset)) }
}

func init() {
	AppendType[string]()
	AppendType[[]byte]()

	AppendType[int8]()
	AppendType[int16]()
	AppendType[int32]()
	AppendType[int64]()

	AppendType[uint8]()
	AppendType[uint16]()
	AppendType[uint32]()
	AppendType[uint64]()

	AppendType[int]()
	AppendType[uint]()

	AppendType[float32]()
	AppendType[float64]()

	AppendType[bool]()

	AppendType[time.Time]()
	AppendType[time.Duration]()
}

// PtrGetter
// return a function that can get this field ptr from an instance ptr.
// calling `AppendType[T any]()` with the type of field, will return a faster function.
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
