package lion

import (
	"fmt"
	"iter"
	"reflect"
	"unsafe"
)

type TypeInfo struct {
	Name   string
	GoType reflect.Type

	fields    []Field
	offsetmap map[int64]*Field

	PtrAny    any
	PtrUnsafe unsafe.Pointer
	PtrNum    int64
}

var (
	ptrs = map[reflect.Type]any{}
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

func ptrof(gotype reflect.Type) any {
	pv, ok := ptrs[gotype]
	if ok {
		return pv
	}
	pv = reflect.New(gotype).Interface()
	ptrs[gotype] = pv
	return pv
}

func Typeof[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

var (
	typeinfocache = map[reflect.Type]*TypeInfo{}
)

func TypeInfoOf[T any]() *TypeInfo {
	gotype := Typeof[T]()
	ti, ok := typeinfocache[gotype]
	if ok {
		return ti
	}
	ti = makeTypeinfo(gotype, Ptr[T]())
	typeinfocache[gotype] = ti
	return ti
}

func makeTypeinfo(gotype reflect.Type, ptr any) *TypeInfo {
	ptrv := reflect.ValueOf(ptr)
	uptr := ptrv.UnsafePointer()

	ti := &TypeInfo{
		Name:      gotype.Name(),
		GoType:    gotype,
		PtrAny:    ptr,
		PtrUnsafe: uptr,
		PtrNum:    int64(uintptr(uptr)),
	}

	walk(&ti.fields, ti.GoType, ptrv, ti.PtrNum)

	if len(ti.fields) > 15 {
		ti.offsetmap = map[int64]*Field{}
		for i := 0; i < len(ti.fields); i++ {
			ptr := &ti.fields[i]
			ti.offsetmap[ptr.offset] = ptr
		}
	}
	for f := range ti.Fields(nil) {
		f.typeinfo = ti
		f._PtrGetter()
		f._Setter()
	}
	return ti
}

func walk(fs *[]Field, gotype reflect.Type, ptrv reflect.Value, begin int64) {
	vv := ptrv.Elem()

	for i := 0; i < gotype.NumField(); i++ {
		sf := gotype.Field(i)
		fv := vv.Field(i)
		fptr := fv.Addr()
		if sf.Anonymous {
			if sf.Type.Kind() != reflect.Struct {
				panic(fmt.Errorf("lion: pointer emmbed struct is not supported. %s.%s", gotype, sf.Type))
			}
			baseoffset := int64(sf.Offset)
			sti := makeTypeinfo(sf.Type, ptrof(sf.Type))

			for idx := range sti.fields {
				stf := &sti.fields[idx]
				f := Field{
					offset: baseoffset + stf.offset,
				}
				if stf.ref != nil {
					f.ref = stf.ref
				} else {
					f.ref = stf
				}
				*fs = append(*fs, f)
			}
			continue
		}
		field := Field{
			field:  sf,
			offset: int64(fptr.Pointer()) - begin,
		}
		*fs = append(*fs, field)
	}
}

func (ti *TypeInfo) FieldByOffset(offset int64) *Field {
	if ti.offsetmap != nil {
		return ti.offsetmap[offset]
	}
	for idx := range ti.fields {
		fp := &ti.fields[idx]
		if fp.offset == offset {
			return fp
		}
	}
	panic(fmt.Errorf("reflectx: bad offset, %s, %d", ti.GoType, offset))
}

func (ti *TypeInfo) FieldByUnsafePtr(ptr unsafe.Pointer) *Field {
	return ti.FieldByOffset(int64(uintptr(ptr)) - ti.PtrNum)
}

func (ti *TypeInfo) FieldByPtr(ptr any) *Field {
	return ti.FieldByUnsafePtr(reflect.ValueOf(ptr).UnsafePointer())
}

type FieldsOptions struct {
	TagName      string
	OnlyExported bool
}

func (ti *TypeInfo) Fields(opts *FieldsOptions) iter.Seq[*Field] {
	if opts == nil {
		opts = &FieldsOptions{}
	}
	if opts.TagName == "" && !opts.OnlyExported {
		return func(yield func(*Field) bool) {
			fc := len(ti.fields)
			for i := 0; i < fc; i++ {
				if !yield(&(ti.fields[i])) {
					break
				}
			}
		}
	}
	return func(yield func(*Field) bool) {
		fc := len(ti.fields)
		for i := 0; i < fc; i++ {
			fp := &(ti.fields[i])
			if opts.TagName != "" {
				tag := fp.Tag(opts.TagName)
				if tag.Name == "-" {
					continue
				}
			}
			if opts.OnlyExported {
				if !fp.StructField().IsExported() {
					continue
				}
			}
			if !yield(&(ti.fields[i])) {
				break
			}
		}
	}
}
