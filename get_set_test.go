package reflectx

import (
	"reflect"
	"testing"
	"unsafe"
)

type A struct {
	A1 string `json:"a1"`
}

func noop(v any) {}

func BenchmarkNormalGetFieldPtrByIndex(b *testing.B) {
	f, _ := reflect.TypeOf(A{}).FieldByName("A1")

	val := A{}

	vv := reflect.ValueOf(&val).Elem()
	for i := 0; i < b.N; i++ {
		fv := vv.FieldByIndex(f.Index).Addr().Interface()
		noop(fv)
	}
}

func BenchmarkGetFieldPtrByOffset(b *testing.B) {
	var ptr = &A{}
	a1offset := int64(uintptr(unsafe.Pointer(&ptr.A1))) - int64(uintptr(unsafe.Pointer(ptr)))

	f, _ := reflect.TypeOf(A{}).FieldByName("A1")

	val := A{}
	valptr := unsafe.Pointer(&val)

	for i := 0; i < b.N; i++ {
		fv := reflect.NewAt(f.Type, unsafe.Add(valptr, a1offset)).Interface()
		noop(fv)
	}
}
