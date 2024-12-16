package reflectx

import (
	"fmt"
	"reflect"
)

func (ti *TypeInfo[M]) Mixed(ptr any, subti *TypeInfo[M]) {
	baseoffset := int64(reflect.ValueOf(ptr).UnsafeAddr()) - ti.PtrNum
	vv := reflect.ValueOf(ti.PtrAny).Elem()

	found := false
	for i := 0; i < ti.GoType.NumField(); i++ {
		sf := ti.GoType.Field(i)
		fptrv := vv.Field(i).Addr().UnsafeAddr()
		foffset := int64(fptrv) - ti.PtrNum
		if sf.Anonymous && sf.Type == subti.GoType && foffset == baseoffset {
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
		*fv = subf
	}
}
