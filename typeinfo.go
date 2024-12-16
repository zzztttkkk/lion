package reflectx

import (
	"reflect"
	"strings"
	"unsafe"
)

type TypeInfo[M any] struct {
	Name   string
	GoType reflect.Type

	Fields    []Field[M]
	offsetmap map[int64]*Field[M]

	PtrAny    any
	PtrUnsafe unsafe.Pointer
	PtrNum    int64
}

var (
	typeinfos = map[reflect.Type]any{}
	ptrs      = map[reflect.Type]any{}
)

func Ptr[T any]() *T {
	gotype := Typeof[T]()
	pv, ok := ptrs[gotype]
	if ok {
		return pv.(*T)
	}
	ptr := new(T)
	ptrs[gotype] = ptr
	return ptr
}

func Typeof[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func TypeInfoOf[T any, M any]() *TypeInfo[M] {
	gotype := Typeof[T]()
	ti, ok := typeinfos[gotype]
	if ok {
		return ti.(*TypeInfo[M])
	}
	reg := RegisterOf[M]()
	obj := makeTypeinfo(reg, gotype, Ptr[T]())
	typeinfos[gotype] = obj
	return obj
}

func makeTypeinfo[M any](reg *_Register[M], gotype reflect.Type, ptr any) *TypeInfo[M] {
	ptrv := reflect.ValueOf(ptr)
	uptr := ptrv.UnsafePointer()
	ti := &TypeInfo[M]{
		GoType:    gotype,
		PtrAny:    ptr,
		PtrUnsafe: uptr,
		PtrNum:    int64(uintptr(uptr)),
	}
	addfields(reg, &ti.Fields, ti.GoType, ptrv, ti.PtrNum)
	if len(ti.Fields) > 15 {
		ti.offsetmap = map[int64]*Field[M]{}
		for i := 0; i < len(ti.Fields); i++ {
			ptr := &ti.Fields[i]
			ti.offsetmap[ptr.Offset] = ptr
		}
	}
	return ti
}

func gettagname(v string) string {
	parts := strings.Split(v, ",")
	return strings.TrimSpace(parts[0])
}

func gettag(sf *reflect.StructField, tags ...string) string {
	for _, tag := range tags {
		v := sf.Tag.Get(tag)
		if v != "" {
			return v
		}
	}
	return ""
}

func addfields[M any](reg *_Register[M], fs *[]Field[M], gotype reflect.Type, ptrv reflect.Value, begin int64) {
	vv := ptrv.Elem()

	for i := 0; i < gotype.NumField(); i++ {
		sf := gotype.Field(i)
		tag := gettag(&sf, reg.tagnames...)
		if tag == "-" {
			continue
		}
		if !reg.unexposed && !sf.IsExported() {
			continue
		}
		fv := vv.Field(i)
		fptr := fv.Addr()
		if sf.Anonymous {
			var _fs []Field[M]
			addfields(reg, &_fs, sf.Type, fptr, begin)
			*fs = append(*fs, _fs...)
			continue
		}
		field := Field[M]{
			Name:   gettagname(tag),
			Field:  sf,
			Offset: int64(fptr.Pointer()) - begin,
		}
		if field.Name == "" {
			field.Name = sf.Name
		}
		*fs = append(*fs, field)
	}
}

func (ti *TypeInfo[M]) FieldByOffset(offset int64) *Field[M] {
	if ti.offsetmap != nil {
		return ti.offsetmap[offset]
	}
	for idx := range ti.Fields {
		fp := &ti.Fields[idx]
		if fp.Offset == offset {
			return fp
		}
	}
	return nil
}

func (ti *TypeInfo[M]) FieldByUnsafePtr(ptr unsafe.Pointer) *Field[M] {
	return ti.FieldByOffset(int64(uintptr(ptr)) - ti.PtrNum)
}

func (ti *TypeInfo[M]) FieldByPtr(ptr any) *Field[M] {
	return ti.FieldByUnsafePtr(reflect.ValueOf(ptr).UnsafePointer())
}
