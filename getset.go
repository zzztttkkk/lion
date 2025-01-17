package lion

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

var (
	ptrgetters = map[reflect.Type]func(insptr unsafe.Pointer, offset int64) any{}
	setters    = map[reflect.Type]func(insptr unsafe.Pointer, offset int64, val any){}
)

// AppendType
// You can noly call this function in `init`, because there is no mutex here.
func AppendType[T any]() {
	// ptr getters
	ptrgetters[Typeof[T]()] = func(insptr unsafe.Pointer, offset int64) any { return (*T)(unsafe.Add(insptr, offset)) }
	ptrgetters[Typeof[*T]()] = func(insptr unsafe.Pointer, offset int64) any { return (**T)(unsafe.Add(insptr, offset)) }
	ptrgetters[Typeof[[]T]()] = func(insptr unsafe.Pointer, offset int64) any { return (*[]T)(unsafe.Add(insptr, offset)) }
	ptrgetters[Typeof[[]*T]()] = func(insptr unsafe.Pointer, offset int64) any { return (*[]*T)(unsafe.Add(insptr, offset)) }
	ptrgetters[Typeof[sql.Null[T]]()] = func(insptr unsafe.Pointer, offset int64) any { return (*sql.Null[T])(unsafe.Add(insptr, offset)) }

	// setters
	setters[Typeof[T]()] = func(insptr unsafe.Pointer, offset int64, val any) {
		fptr := (*T)(unsafe.Add(insptr, offset))
		*fptr = (val.(T))
	}
	setters[Typeof[*T]()] = func(insptr unsafe.Pointer, offset int64, val any) {
		fptr := (**T)(unsafe.Add(insptr, offset))
		*fptr = (val.(*T))
	}
	setters[Typeof[[]T]()] = func(insptr unsafe.Pointer, offset int64, val any) {
		fptr := (*[]T)(unsafe.Add(insptr, offset))
		*fptr = (val.([]T))
	}
	setters[Typeof[[]*T]()] = func(insptr unsafe.Pointer, offset int64, val any) {
		fptr := (*[]*T)(unsafe.Add(insptr, offset))
		*fptr = (val.([]*T))
	}
	setters[Typeof[sql.Null[T]]()] = func(insptr unsafe.Pointer, offset int64, val any) {
		fptr := (*sql.Null[T])(unsafe.Add(insptr, offset))
		*fptr = (val.(sql.Null[T]))
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

type anyface struct {
	typeuptr unsafe.Pointer
	valuptr  unsafe.Pointer
}

func pack(tuptr unsafe.Pointer, value unsafe.Pointer) any {
	var iv any
	ivptr := (*anyface)(unsafe.Pointer(&iv))
	ivptr.typeuptr = tuptr
	ivptr.valuptr = value
	return iv
}

func (field *Field) _PtrGetter() FieldPtrGetter {
	if field.ptrgetter == nil {
		sf := field.StructField()
		getter, ok := ptrgetters[sf.Type]
		if ok {
			field.ptrgetter = func(insptr unsafe.Pointer) any { return getter(insptr, field.offset) }
		} else {
			ptrtype := reflect.PointerTo(sf.Type)
			ptrtypeuptr := reflect.ValueOf(ptrtype).UnsafePointer()
			field.ptrgetter = func(insptr unsafe.Pointer) any {
				return pack(ptrtypeuptr, unsafe.Add(insptr, field.offset))
			}
		}
	}
	return field.ptrgetter
}

func (field *Field) PtrOf(insptr unsafe.Pointer) any { return field.ptrgetter(insptr) }

func (field *Field) UnsafePtrOf(insptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Add(insptr, field.offset)
}

func memcopy(dst unsafe.Pointer, src unsafe.Pointer, bytes int) {
	type SliceHeader struct {
		Data uintptr
		Len  int
		Cap  int
	}
	var dstsh = SliceHeader{
		Data: uintptr(dst),
		Len:  bytes,
		Cap:  bytes,
	}
	var srcsh = SliceHeader{
		Data: uintptr(src),
		Len:  bytes,
		Cap:  bytes,
	}
	var dstsptr = (*[]byte)(unsafe.Pointer(&dstsh))
	var srcsptr = (*[]byte)(unsafe.Pointer(&srcsh))
	copy(*dstsptr, *srcsptr)
}

var (
	ErrNil = errors.New("lion: src is nil")
)

func (field *Field) _Setter() func(insptr unsafe.Pointer, val any) {
	if field.setter == nil {
		sf := field.StructField()
		setter, ok := setters[sf.Type]
		if ok {
			field.setter = func(insptr unsafe.Pointer, val any) {
				setter(insptr, field.offset, val)
			}
		} else {
			typeuptr := reflect.ValueOf(sf.Type).UnsafePointer()
			size := int(sf.Type.Size())
			isptr := sf.Type.Kind() == reflect.Pointer

			if isptr {
				field.setter = func(insptr unsafe.Pointer, src any) {
					srcface := (*anyface)(unsafe.Pointer(&src))
					if srcface.typeuptr != typeuptr {
						panic(fmt.Errorf("lion: `%v` is not type `%s`", src, sf.Type))
					}
					if srcface.valuptr == nil {
						panic(ErrNil)
					}
					memcopy(unsafe.Add(insptr, field.offset), unsafe.Pointer(&srcface.valuptr), size)
				}
			} else {
				field.setter = func(insptr unsafe.Pointer, src any) {
					srcface := (*anyface)(unsafe.Pointer(&src))
					if srcface.typeuptr != typeuptr {
						panic(fmt.Errorf("lion: `%v` is not type `%s`", src, sf.Type))
					}
					if srcface.valuptr == nil {
						panic(ErrNil)
					}
					memcopy(unsafe.Add(insptr, field.offset), srcface.valuptr, size)
				}
			}
		}
	}
	return field.setter
}

func (field *Field) AssignTo(insptr unsafe.Pointer, val any) {
	field.setter(insptr, val)
}
