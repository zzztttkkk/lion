package reflectx

import (
	"fmt"
	"reflect"
	"unsafe"
)

func (ti *TypeInfo[M]) MixByUnsafePtr(ptr unsafe.Pointer, subti *TypeInfo[M]) *TypeInfo[M] {
	baseoffset := int64(uintptr(ptr)) - ti.PtrNum

	found := false
	for i := 0; i < ti.GoType.NumField(); i++ {
		sf := ti.GoType.Field(i)
		if sf.Anonymous && sf.Type == subti.GoType && int64(sf.Offset) == baseoffset {
			found = true
			break
		}
	}
	if !found {
		panic(fmt.Errorf("reflectx: not found"))
	}
	for _, subf := range subti.Fields {
		coffset := baseoffset + subf.Offset
		fv := ti.FieldByOffset(coffset)
		fv.Name = subf.Name
		fv.Meta = subf.Meta
	}
	return ti
}

func (ti *TypeInfo[M]) Mix(ptr any, subti *TypeInfo[M]) *TypeInfo[M] {
	return ti.MixByUnsafePtr(reflect.ValueOf(ptr).UnsafePointer(), subti)
}
