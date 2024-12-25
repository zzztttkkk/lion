package lion

import (
	"reflect"
)

func FindTypeinfo[M any](fptr any) *TypeInfo[M] {
	reg := RegisterOf[M]()
	ufptr := reflect.ValueOf(fptr).UnsafePointer()
	trygetfield := func(ti *TypeInfo[M]) *Field[M] {
		var f *Field[M]
		defer func() {
			recover()
			f = nil
		}()
		f = ti.FieldByUnsafePtr(ufptr)
		return f
	}
	for gotype := range ptrs {
		typeinfo, ok := reg.typeinfos[gotype]
		if !ok {
			continue
		}
		f := trygetfield(typeinfo)
		if f != nil {
			return typeinfo
		}
	}
	return nil
}
