package lion

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

var (
	preparedTypePtrGetters   = map[reflect.Type]func(insptr unsafe.Pointer, offset int64) any{}
	preparedTypeValueGetters = map[reflect.Type]func(insptr unsafe.Pointer, offset int64) any{}
	preparedTypeSetters      = map[reflect.Type]func(insptr unsafe.Pointer, offset int64, val any){}
)

// AppendType
// You can noly call this function in `init`, because there is no mutex here.
func AppendType[T any]() {
	// ptr getters
	preparedTypePtrGetters[Typeof[T]()] = func(insptr unsafe.Pointer, offset int64) any { return (*T)(unsafe.Add(insptr, offset)) }
	preparedTypePtrGetters[Typeof[*T]()] = func(insptr unsafe.Pointer, offset int64) any { return (**T)(unsafe.Add(insptr, offset)) }
	preparedTypePtrGetters[Typeof[[]T]()] = func(insptr unsafe.Pointer, offset int64) any { return (*[]T)(unsafe.Add(insptr, offset)) }
	preparedTypePtrGetters[Typeof[[]*T]()] = func(insptr unsafe.Pointer, offset int64) any { return (*[]*T)(unsafe.Add(insptr, offset)) }
	preparedTypePtrGetters[Typeof[sql.Null[T]]()] = func(insptr unsafe.Pointer, offset int64) any { return (*sql.Null[T])(unsafe.Add(insptr, offset)) }

	// value getters
	preparedTypeValueGetters[Typeof[T]()] = func(insptr unsafe.Pointer, offset int64) any { return *(*T)(unsafe.Add(insptr, offset)) }
	preparedTypeValueGetters[Typeof[*T]()] = func(insptr unsafe.Pointer, offset int64) any { return *(**T)(unsafe.Add(insptr, offset)) }
	preparedTypeValueGetters[Typeof[[]T]()] = func(insptr unsafe.Pointer, offset int64) any { return *(*[]T)(unsafe.Add(insptr, offset)) }
	preparedTypeValueGetters[Typeof[[]*T]()] = func(insptr unsafe.Pointer, offset int64) any { return *(*[]*T)(unsafe.Add(insptr, offset)) }
	preparedTypeValueGetters[Typeof[sql.Null[T]]()] = func(insptr unsafe.Pointer, offset int64) any { return *(*sql.Null[T])(unsafe.Add(insptr, offset)) }

	// setters
	preparedTypeSetters[Typeof[T]()] = func(insptr unsafe.Pointer, offset int64, val any) {
		fptr := (*T)(unsafe.Add(insptr, offset))
		*fptr = *(anytotype[T](val))
	}
	preparedTypeSetters[Typeof[*T]()] = func(insptr unsafe.Pointer, offset int64, val any) {
		fptr := (**T)(unsafe.Add(insptr, offset))
		*fptr = *(anytotype[*T](val))
	}
	preparedTypeSetters[Typeof[[]T]()] = func(insptr unsafe.Pointer, offset int64, val any) {
		fptr := (*[]T)(unsafe.Add(insptr, offset))
		*fptr = *(anytotype[[]T](val))
	}
	preparedTypeSetters[Typeof[[]*T]()] = func(insptr unsafe.Pointer, offset int64, val any) {
		fptr := (*[]*T)(unsafe.Add(insptr, offset))
		*fptr = *(anytotype[[]*T](val))
	}
	preparedTypeSetters[Typeof[sql.Null[T]]()] = func(insptr unsafe.Pointer, offset int64, val any) {
		fptr := (*sql.Null[T])(unsafe.Add(insptr, offset))
		*fptr = *(anytotype[sql.Null[T]](val))
	}
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

var (
	warnningedTypes = map[reflect.Type]EmptyMeta{}
)

func warnningForType(gotype reflect.Type) {
	_, ok := warnningedTypes[gotype]
	if ok {
		return
	}
	warnningedTypes[gotype] = EmptyMeta{}
	fmt.Printf("reflectx.warnning: reflectx.AppendType[%s]() is not called, please call it to improve performance\r\n", gotype)
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
			warnningForType(sf.Type)
			field.ptrgetter = func(insptr unsafe.Pointer) any {
				return reflect.NewAt(sf.Type, unsafe.Add(insptr, field.offset)).Interface()
			}
		}
	}
	return field.ptrgetter
}

func (field *Field[M]) PtrOfInstance(insptr unsafe.Pointer) any { return field.PtrGetter()(insptr) }

func (field *Field[M]) Getter() _FieldPtrGetter {
	if field.getter == nil {
		sf := field.StructField()
		getter, ok := preparedTypeValueGetters[sf.Type]
		if ok {
			field.getter = func(insptr unsafe.Pointer) any { return getter(insptr, field.offset) }
		} else {
			warnningForType(sf.Type)
			field.getter = func(insptr unsafe.Pointer) any {
				return reflect.NewAt(sf.Type, unsafe.Add(insptr, field.offset)).Elem().Interface()
			}
		}
	}
	return field.getter
}

func (field *Field[M]) ValueOfInstance(insptr unsafe.Pointer) any { return field.Getter()(insptr) }

func (field *Field[M]) Setter() func(insptr unsafe.Pointer, val any) {
	if field.setter == nil {
		sf := field.StructField()
		setter, ok := preparedTypeSetters[sf.Type]
		if ok {
			field.setter = func(insptr unsafe.Pointer, val any) {
				setter(insptr, field.offset, val)
			}
		} else {
			warnningForType(sf.Type)
			field.setter = func(insptr unsafe.Pointer, val any) {
				reflect.NewAt(sf.Type, unsafe.Add(insptr, field.offset)).Elem().Set(reflect.ValueOf(val))
			}
		}
	}
	return field.setter
}

func (field *Field[M]) ChangeInstance(insptr unsafe.Pointer, val any) {
	field.Setter()(insptr, val)
}
