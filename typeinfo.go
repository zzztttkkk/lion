package reflectx

import (
	"reflect"
	"strings"
	"unsafe"
)

type TypeInfo struct {
	Name   string
	GoType reflect.Type

	Fields    []Field
	offsetmap map[int64]*Field

	Ptr    unsafe.Pointer
	PtrNum int64
}

var (
	typeinfos = map[reflect.Type]*TypeInfo{}
)

func typeof[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func TypeInfoOf[T any]() *TypeInfo {
	gotype := typeof[T]()
	ti, ok := typeinfos[gotype]
	if ok {
		return ti
	}

	obj := &TypeInfo{
		GoType: gotype,
	}
	return obj
}

func makeTypeinfo(tagnames []string, gotype reflect.Type, ptr any) *TypeInfo {
	ptrv := reflect.ValueOf(ptr)
	uptr := ptrv.UnsafePointer()
	ti := &TypeInfo{
		GoType: gotype,
		Ptr:    uptr,
		PtrNum: int64(uintptr(uptr)),
	}
	addfields(&ti.Fields, ti.GoType, tagnames, ptrv, ti.PtrNum)
	if len(ti.Fields) > 12 {
		ti.offsetmap = map[int64]*Field{}
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

func addfields(fs *[]Field, gotype reflect.Type, tagnames []string, ptrv reflect.Value, begin int64) {
	vv := ptrv.Elem()

	for i := 0; i < gotype.NumField(); i++ {
		sf := gotype.Field(i)
		tag := gettag(&sf, tagnames...)
		if !sf.IsExported() || tag == "-" {
			continue
		}
		fv := vv.Field(i)
		fptr := fv.Addr()
		if sf.Anonymous {
			var _fs []Field
			addfields(&_fs, sf.Type, tagnames, fptr, begin)
			*fs = append(*fs, _fs...)
			return
		}
		field := Field{
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

func (ti *TypeInfo) fieldByOffset(offset int64) *Field {
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

func (ti *TypeInfo) FieldByUnsafePtr(ptr unsafe.Pointer) *Field {
	return ti.fieldByOffset(int64(uintptr(ptr)) - ti.PtrNum)
}

func (ti *TypeInfo) FieldByPtr(ptr any) *Field {
	return ti.FieldByUnsafePtr(reflect.ValueOf(ptr).UnsafePointer())
}
